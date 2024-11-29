package mainform

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms"
	"faxsender/src/ui/forms/mainform/tabs"
	"faxsender/src/utilities"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

// MainForm represents the main user interface of the application.
type MainForm struct {
	app    *fyne.App
	window *fyne.Window

	tabManagement *tabs.TabManagement

	apiUI api.IApiUICalls

	forms.IBaseForm
}

// NewMainForm creates a new instance of the MainForm.
//
// Steps:
// 1. Create a new Fyne application.
// 2. Create a new Fyne window.
// 3. Initialize application settings and window size.
// 4. Create a new API user interface (ApiUI) instance with the server port.
//
// Parameters:
//
//	None
//
// Returns:
//   - *MainForm: A pointer to the created MainForm instance.
func NewMainForm() *MainForm {
	app := app.New()
	window := app.NewWindow(utilities.APP_NAME)

	forms.InitApp(&app)
	forms.InitForms(&window)
	forms.SetSize(&window, 800, 450)

	apiUI := api.NewApiServerDirectCalls()

	return &MainForm{
		app:    &app,
		window: &window,
		apiUI:  apiUI,
	}
}

// InitControls initializes the controls of the MainForm.
//
// Steps:
// 1. Create a new TabManagement instance with the API user interface and window.
// 2. Set the window content to the TabManagement instance.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *MainForm) InitControls() {
	f.tabManagement = tabs.NewTabeManagement(&f.apiUI, f.window)
	(*f.window).SetContent(f.tabManagement.GetTabContainer())
}

// Show displays the MainForm and runs the Fyne application.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *MainForm) Show() {
	(*f.window).ShowAndRun()
}
