package tabs

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms/sendfaxform"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
)

// SendFaxTab represents a tab for sending a fax.
type SendFaxTab struct {
	api    *api.IApiUICalls
	parent *fyne.Window

	tabItem *container.TabItem

	sendFaxForm *sendfaxform.SendFaxForm

	signalFunc TabManagementSignal

	ITab
}

const (
	DEFAULT_PATH = ""
)

// NewSendFaxTab creates a new SendFaxTab instance.
// Parameters:
//   - api: An instance of the API client for UI calls.
//   - parent: The parent Fyne window.
//
// Returns:
//   - A pointer to the newly created SendFaxTab.
func NewSendFaxTab(api *api.IApiUICalls, parent *fyne.Window) *SendFaxTab {
	return &SendFaxTab{
		api:    api,
		parent: parent,
	}
}

// initUI initializes the user interface components of the tab.
func (s *SendFaxTab) initUI() {
	s.sendFaxForm = sendfaxform.NewSendFaxForm(DEFAULT_PATH, s.parent)
	s.sendFaxForm.InitControls()
	s.sendFaxForm.SignalFunc = s.sendFaxFormSignal
}

// loadData loads data for the tab.
func (s *SendFaxTab) loadData() {

}

// GetTab returns the TabItem associated with this tab.
// Returns:
//   - The TabItem for the SendFaxTab.
func (s *SendFaxTab) GetTab() *container.TabItem {
	if s.tabItem == nil {
		s.tabItem = container.NewTabItem("Send Fax", s.sendFaxForm.GetMainContainer())
		s.tabItem.Icon = theme.MailSendIcon()
	}
	return s.tabItem
}

// setSignalFunc sets the signal function for the tab.
// Parameters:
//   - signalFunc: The function to handle signals related to tab management.
func (s *SendFaxTab) setSignalFunc(signalFunc TabManagementSignal) {
	s.signalFunc = signalFunc
}

// sendFaxFormSignal is the signal handler for the sendFaxForm.
// Parameters:
//   - args: A list of integers representing signal arguments.
func (s *SendFaxTab) sendFaxFormSignal(args ...int) {
	s.signalFunc(SIGNAL_FAX_SENT_SUCCESS)
}
