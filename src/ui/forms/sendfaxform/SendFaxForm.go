package sendfaxform

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms"
	"faxsender/src/utilities"
	"faxsender/src/utilities/logger"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// Constants for retry combo box options.
const (
	RETRYCOMBO_DEFAULT_STRING string = "Choose a number to retry"
	RETRYCOMBO_ONE            string = "1"
	RETRYCOMBO_TWO            string = "2"
	RETRYCOMBO_THREE          string = "3"
	RETRYCOMBO_FOUR           string = "4"
	RETRYCOMBO_FIVE           string = "5"

	PHONE_LIST_DEFAULT_STRING string = "Choose Caller ID*"

	IS_NOT_SELECTED_STRING string = ""
)

type SendFaxFormSignal func(...int)

// SendFaxForm is a form for sending a fax, providing user input fields and options.
type SendFaxForm struct {
	app    *fyne.App
	window *fyne.Window

	firstNameEntry   *widget.Entry
	lastNameEntry    *widget.Entry
	emailEntry       *widget.Entry
	faxNumberEntry   *widget.Entry
	companyEntry     *widget.Entry
	descriptionEntry *widget.Entry
	custom1Entry     *widget.Entry
	custom2Entry     *widget.Entry
	custom3Entry     *widget.Entry
	titleEntry       *widget.Entry
	retryEntry       *widget.Select
	accountPhoneList *widget.Select
	sendButton       *widget.Button
	selectContainer  container.Scroll
	filePathLable    *widget.Entry

	infoEntryLayout *fyne.Container

	formLayout        *fyne.Container
	coverPageCheckbox *fyne.Container
	printCheckbox     *fyne.Container
	buttonContainer   *fyne.Container

	apiUI api.IApiUICalls

	forms.IBaseForm

	accountResponse []api.AccountResponse
	allAccounts     []api.AccountResponse
	transmission    api.Transmission
	contact         api.Contact
	documentRecord  api.DocumentRecord
	fileModel       api.SendFileInfo

	filePath     string
	fileContents []byte
	SignalFunc   SendFaxFormSignal
}

// NewSendFaxForm creates a new instance of SendFaxForm.
//
// Parameters:
//   - filePath: The file path of the document to be faxed.
//
// Returns:
//   - *SendFaxForm: The created SendFaxForm instance.
func NewSendFaxForm(filePath string, parent *fyne.Window) *SendFaxForm {
	app := app.New()
	window := app.NewWindow(utilities.APP_NAME)

	forms.InitApp(&app)
	forms.InitForms(&window)
	forms.SetSize(&window, 700, 400)

	apiUI := api.NewApiServerDirectCalls()

	if parent != nil {
		window = *parent
	}

	return &SendFaxForm{
		app:    &app,
		window: &window,
		apiUI:  apiUI,

		filePath: utilities.CleanupFilePath(filePath),
		fileModel: api.SendFileInfo{
			ContentType: utilities.EMPTY_FILE_EXTENSION,
		},
	}
}

func (r *SendFaxForm) CheckPathWithPanic() {
	if utilities.ExtractFileExtension(r.filePath) == utilities.EMPTY_FILE_EXTENSION {
		msg := fmt.Sprintf("the file path '%s' is not a valid file extension", r.filePath)
		logger.Inst().Error(msg)
		utilities.ExecuteOnTerminal("zenity", "--error", "title", "error in path", "--text", msg)
		os.Exit(utilities.ERROR_CDOE_INVALID_FILE_EXTENSION)
	}
}

// InitControls initializes the UI controls of the SendFaxForm.
//
// Steps:
// 1. Initialize input entries (text fields).
// 2. Initialize information layout.
// 3. Initialize checkboxes.
// 4. Initialize retry combo box.
// 5. Initialize phone list combo box.
// 6. Initialize send button.
// 7. Initialize form layout.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) InitControls() {

	recipientInfoTitle := widget.NewLabel("Recipient Information")
	uploadTitle := widget.NewLabel("Upload Your Fax Document*")
	uploadTitle.TextStyle = fyne.TextStyle{Bold: true}
	fileLable := widget.NewEntry()
	fileLable.Text = "File Path:"
	fileLable.Disable()

	f.filePathLable = widget.NewEntry()
	f.filePathLable.Text = f.filePath
	f.filePathLable.Disable()

	browseButton := widget.NewButton("", func() {
		f.openFileDialog()
	})
	browseButton.Icon = theme.FileIcon()
	browseButton.SetText("Upload")

	fileContainer := container.NewBorder(nil, nil, fileLable, browseButton, f.filePathLable)

	f.initInputEntries()
	f.initInformationsLayout()
	f.initCheckBoxes()
	f.initRetryCombobox()
	f.initPhoneListCombobox()
	f.initSendButton()
	f.initFormLayout(uploadTitle, recipientInfoTitle, fileContainer)

	(*f.window).SetContent(f.formLayout)
	f.InitAccountPhoneListOptions()

}

// initInformationsLayout initializes the layout for recipient information.
//
// Steps:
// 1. Create text entry fields for various recipient information.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initInformationsLayout() {
	f.infoEntryLayout = container.NewVBox(
		container.NewGridWithColumns(3,
			f.firstNameEntry,
			f.lastNameEntry,
			f.emailEntry,
		),
		container.NewGridWithColumns(3,
			f.faxNumberEntry,
			f.companyEntry,
		),
		container.NewGridWithColumns(1,
			f.descriptionEntry,
		),
		container.NewGridWithColumns(3,
			f.custom1Entry,
			f.custom2Entry,
			f.custom3Entry,
		),
	)
}

// initFormLayout initializes the overall layout of the form.
//
// Steps:
// 1. Create the layout with multiple sections using container layouts.
//
// Parameters:
//   - recipientInfoTitle: Label for recipient information.
//   - fileContainer: Container holding file-related elements//
//
// Returns:
//
//	None
func (f *SendFaxForm) initFormLayout(uploadTitle *widget.Label, recipientInfoTitle *widget.Label, fileContainer *fyne.Container) {
	f.formLayout = container.NewVBox(
		f.titleEntry,
		uploadTitle,
		container.NewVBox(fileContainer, recipientInfoTitle),
		f.infoEntryLayout,
		container.NewGridWithColumns(5, f.coverPageCheckbox, f.printCheckbox),
		container.NewGridWithColumns(3, f.retryEntry, f.accountPhoneList, f.sendButton),
	)
}

// initSendButton initializes the send button.
//
// Steps:
// 1. Create a "Send" button with an icon.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initSendButton() {
	f.sendButton = widget.NewButton("Send", f.onSendClick)
	f.sendButton.Icon = theme.MailSendIcon()
	f.sendButton.Resize(fyne.NewSize(200, 50))
}

// initPhoneListCombobox initializes the combo box for selecting a phone number.
//
// Steps:
// 1. Create a select widget for phone numbers.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initPhoneListCombobox() {
	f.accountPhoneList = widget.NewSelect([]string{PHONE_LIST_DEFAULT_STRING}, func(selected string) {
		if selected != IS_NOT_SELECTED_STRING {
			f.handleAccountSelection(selected)
		}
	})

	f.accountPhoneList.PlaceHolder = PHONE_LIST_DEFAULT_STRING
}

// initRetryCombobox initializes the combo box for selecting retry options.
//
// Steps:
// 1. Create a select widget for retry options.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initRetryCombobox() {
	f.retryEntry = widget.NewSelect([]string{RETRYCOMBO_DEFAULT_STRING,
		RETRYCOMBO_ONE,
		RETRYCOMBO_TWO,
		RETRYCOMBO_THREE,
		RETRYCOMBO_FOUR,
		RETRYCOMBO_FIVE},
		func(selected string) {

			if selected != RETRYCOMBO_DEFAULT_STRING {
				f.handleRetry(selected)
			}
		})
	f.retryEntry.PlaceHolder = RETRYCOMBO_DEFAULT_STRING
}

// initInputEntries initializes the input fields for various user information.
//
// Steps:
// 1. Create text entry fields for different user information.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initInputEntries() {
	f.firstNameEntry = widget.NewEntry()
	f.firstNameEntry.PlaceHolder = "First Name"

	f.lastNameEntry = widget.NewEntry()
	f.lastNameEntry.PlaceHolder = "Last Name"

	f.emailEntry = widget.NewEntry()
	f.emailEntry.PlaceHolder = "Email"

	f.faxNumberEntry = widget.NewEntry()
	f.faxNumberEntry.PlaceHolder = "Fax Number*"

	f.companyEntry = widget.NewEntry()
	f.companyEntry.PlaceHolder = "Company Name"

	f.descriptionEntry = widget.NewEntry()
	f.descriptionEntry.PlaceHolder = "Description"

	f.custom1Entry = widget.NewEntry()
	f.custom1Entry.PlaceHolder = "Custom 1"

	f.custom2Entry = widget.NewEntry()
	f.custom2Entry.PlaceHolder = "Custom 2"

	f.custom3Entry = widget.NewEntry()
	f.custom3Entry.PlaceHolder = "Custom 2"

	f.titleEntry = widget.NewEntry()
	f.titleEntry.PlaceHolder = "Title*"
}

// initCheckBoxes initializes checkboxes for cover page and print options.
//
// Steps:
// 1. Create checkboxes for cover page and print.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) initCheckBoxes() {
	f.coverPageCheckbox = container.NewHBox(widget.NewCheck("Cover Page", func(checked bool) {
		if checked {
			f.transmission.IsCoverPage = utilities.WITH_COVER
		} else {
			f.transmission.IsCoverPage = utilities.WITHOUT_COVER
		}
	}))

	f.printCheckbox = container.NewHBox(widget.NewCheck("Print", func(checked bool) {
		if checked {
			f.transmission.IsPrint = utilities.WITH_PRINT
		} else {
			f.transmission.IsPrint = utilities.WITHOUT_PRINT
		}
	}))
	f.transmission.IsPrint = utilities.WITH_PRINT

	f.printCheckbox.Hidden = true
}

// InitAccountPhoneListOptions initializes the phone list combo box with account phone numbers.
//
// Steps:
// 1. Fetch all accounts from the API.
// 2. Populate the phone list options with account phone numbers.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) InitAccountPhoneListOptions() {
	var phoneNumbers []string

	accounts, err := f.apiUI.GetAllAccounts()
	f.allAccounts = accounts
	if err != nil {
		logger.Inst().Error(err.Error())
		return
	}
	for _, account := range f.allAccounts {
		phoneNumbers = append(phoneNumbers, account.Phone)
	}
	f.accountPhoneList.Options = append([]string{PHONE_LIST_DEFAULT_STRING}, phoneNumbers...)
}

// handleAccountSelection handles the selection of an account phone number.
//
// Parameters:
//   - selectedPhone: The selected phone number from the combo box.
//
// Returns:
//   - accountID: The account ID associated with the selected phone number.
func (f *SendFaxForm) handleAccountSelection(selectedPhone string) (accountID string) {
	if selectedPhone != PHONE_LIST_DEFAULT_STRING {
		for _, account := range f.allAccounts {
			if account.Phone == selectedPhone {
				f.transmission.AccountID = account.AccountID
				println("Selected Phone:", selectedPhone, "Account ID:", account.AccountID)
				break
			}
		}
	}
	return f.transmission.AccountID
}

// handleRetry handles the selection of retry options.
//
// Parameters:
//   - selectedRetry: The selected retry option from the combo box.
//
// Returns:
//   - retryData: The retry data associated with the selected option.
func (f *SendFaxForm) handleRetry(selectedRetry string) (retryData string) {
	if selectedRetry != RETRYCOMBO_DEFAULT_STRING {
		f.transmission.TryAllowed = selectedRetry
	}
	println("Selected Retry time:", f.transmission.TryAllowed)
	return f.transmission.TryAllowed
}

// readFileContents reads the contents of the file specified by the file path.
//
// Steps:
// 1. Check if the file extension is valid.fileAbsolutePath
// Parameters:
//   - filePath: The path of the file to be read.
//
// Returns:
//   - []byte: The contents of the file.
func (f *SendFaxForm) readFileContents(filePath string) []byte {

	if utilities.ExtractFileExtension(filePath) == utilities.EMPTY_FILE_EXTENSION {
		logger.Inst().Error("Invalid file extension")
		forms.ShowError("Invalid file extension", f.window)
		return nil
	}

	f.fileContents, _ = os.ReadFile(filePath)

	return f.fileContents
}

func (f *SendFaxForm) openFileDialog() {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			logger.Inst().Error(err.Error())
			forms.ShowError("Failed to open file dialog", f.window)
			return
		}

		if reader == nil {

			logger.Inst().Info("Failed to init reader")
			return
		}
		defer reader.Close()

		fileUrl := reader.URI()
		fileAbsolutePath := strings.ReplaceAll(fileUrl.String(), "file://", "")
		f.filePathLable.SetText(fileAbsolutePath)
		f.filePath = fileAbsolutePath
	}, *f.window)
	fileDialog.SetFilter(storage.NewExtensionFileFilter(utilities.AllValidExtensions()))
	fileDialog.Show()
}

// onSendClick handles the click event of the "Send" button.
//
// Steps:
// 1. Prepare contact, document record, transmission, and file model data.
// 2. Call the API to send the fax with the prepared data.
// 3. Handle any errors and display an error message if necessary.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) onSendClick() {

	if f.faxNumberEntry.Text == "" || f.titleEntry.Text == "" || f.transmission.AccountID == "" {
		forms.ShowError("Required fields cannot be empty", f.window)
		return
	}

	f.sendButton.Disable()

	f.contact = api.Contact{
		FirstName: f.firstNameEntry.Text,
		LastName:  f.lastNameEntry.Text,
		Email:     f.emailEntry.Text,
		Phone:     f.faxNumberEntry.Text,
	}
	f.documentRecord = api.DocumentRecord{
		Title:       f.titleEntry.Text,
		Description: f.descriptionEntry.Text,
	}

	f.transmission = api.Transmission{
		Title:       f.titleEntry.Text,
		AccountID:   f.transmission.AccountID,
		IsCoverPage: f.transmission.IsCoverPage,
		IsPrint:     f.transmission.IsPrint,
		TryAllowed:  f.transmission.TryAllowed,
	}

	f.fileModel = api.SendFileInfo{
		ContentType: f.fileModel.ContentType,
	}

	f.readFileContents(f.filePath)
	fileExtension := utilities.ExtractFileExtension(f.filePath)
	f.fileModel.ContentType = utilities.GetContentType(fileExtension)

	err := f.apiUI.SendFax(f.contact, f.documentRecord, f.transmission, f.fileContents, f.fileModel)

	f.sendButton.Enable()

	if err != nil {
		logger.Inst().Error(err.Error())
		forms.ShowError("error in sending the fax", f.window)
		return
	}
	forms.ShowInfo("success", "the fax has been sent!", f.window)
	f.SignalFunc()
}

func (f *SendFaxForm) GetMainContainer() *fyne.Container {
	return f.formLayout
}

// Show displays and runs the SendFaxForm.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *SendFaxForm) Show() {
	(*f.window).ShowAndRun()
}
