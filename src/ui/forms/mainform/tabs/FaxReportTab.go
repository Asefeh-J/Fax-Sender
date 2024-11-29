package tabs

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const NUMBER_OF_FAXES = 10

// FaxReportTab is a tab for displaying outbound fax reports.
type FaxReportTab struct {
	api    *api.IApiUICalls
	parent *fyne.Window

	tabItem       *container.TabItem
	mainContainer *fyne.Container

	headerContainer *fyne.Container
	faxList         *fyne.Container
	updateButton    *widget.Button

	ITab
	lastFaxes []api.FaxData
}

// NewFaxReportTab creates a new instance of FaxReportTab.
//
// Parameters:
//   - apiInst: An instance of the ApiUI for API interactions.
//   - parent: The parent window associated with the tab.
//
// Returns:
//   - *FaxReportTab: The created FaxReportTab instance.
func NewFaxReportTab(apiInst *api.IApiUICalls, parent *fyne.Window) *FaxReportTab {
	return &FaxReportTab{
		api:    apiInst,
		parent: parent,
	}
}

// initUI initializes the UI components of the FaxReportTab.
//
// Steps:
// 1. Create and configure UI components such as labels, buttons, and containers.
// 2. Set up the layout of the components using containers and layouts.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *FaxReportTab) initUI() {

	f.headerContainer = fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(5),
		widget.NewLabel("Date and Time"),
		widget.NewLabel("Title"),
		widget.NewLabel("Destination Fax"),
		widget.NewLabel("Caller ID"),
		widget.NewLabel("Status"),
	)

	f.updateButton = widget.NewButton("Update", f.onLoadClick)
	f.updateButton.Icon = theme.DownloadIcon()

	f.faxList = container.NewVBox()

	f.mainContainer = container.NewGridWithRows(1,
		container.NewVScroll(container.NewVBox(container.NewHBox(f.updateButton),
			f.headerContainer,
			f.faxList,
		)),
	)
}

// loadData fetches the latest fax data from the API.
//
// Steps:
// 1. Call the API to retrieve the latest fax data.
// 2. Handle any errors during the data retrieval process.
//
// Returns:
//   - bool: True if data is loaded successfully, false otherwise.
func (f *FaxReportTab) loadData() bool {
	Faxes, err := (*f.api).GetLastFaxes(NUMBER_OF_FAXES)
	f.lastFaxes = Faxes
	if err != nil {
		fmt.Println("Error fetching data:", err)
		// Handle error (log, display a message, etc.)
		return false
	}

	return true
}

// loadDataIntoUI populates the UI components with the fetched fax data.
//
// Parameters:
//   - lastFaxes: A slice containing the latest fax data.
//
// Returns:
//
//	None
func (f *FaxReportTab) loadDataIntoUI(lastFaxes []api.FaxData) {
	f.faxList.Objects = nil

	for _, fax := range lastFaxes {
		labels := []fyne.CanvasObject{
			widget.NewLabel(fax.DateTime.Format("2006_01_02 15:04:05")),
			widget.NewLabel(fax.Title),
			widget.NewLabel(fax.DestinationFax),
			widget.NewLabel(fax.CallerID),
			widget.NewLabel(fax.Status),
		}
		rowContainer := fyne.NewContainerWithLayout(layout.NewGridLayoutWithColumns(5), labels...)
		f.faxList.Add(rowContainer)
	}
}

// GetTab returns the TabItem associated with the FaxReportTab.
//
// Steps:
// 1. Return the TabItem associated with the tab.
//
// Returns:
//   - *container.TabItem: The TabItem associated with the tab.
func (f *FaxReportTab) GetTab() *container.TabItem {
	if f.tabItem == nil {
		f.tabItem = container.NewTabItem("Outbound Fax List", f.mainContainer)
		f.tabItem.Icon = theme.ComputerIcon()
	}
	return f.tabItem
}

// IsDataLoaded checks if the fax data is loaded.
//
// Returns:
//   - bool: True if the data is loaded, false otherwise.
func (f *FaxReportTab) IsDataLoaded() bool {
	return f.loadData()
}

// onLoadClick is the callback function for the update button.
//
// Steps:
// 1. Fetch the latest fax data from the API.
// 2. Populate the UI components with the fetched data.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (f *FaxReportTab) onLoadClick() {
	f.faxList.Hide()
	if !f.loadData() {
		forms.ShowError("data cannot be loaded!", f.parent)
		return
	}

	f.loadDataIntoUI(f.lastFaxes)
	f.faxList.Show()
	forms.ShowInfo("Data Loaded", "the data is loaded!", f.parent)
}
