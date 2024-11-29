package tabs

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms"
	"faxsender/src/utilities/logger"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// Constants for tab management signals.
const (
	SIGNAL_SETTINGS_SAVED   = 1
	SIGNAL_LOGOUT           = 2
	SIGNAL_FAX_SENT_SUCCESS = 3
)

// TabManagement represents the management of different tabs in the UI.
type TabManagement struct {
	tabs *widget.TabContainer

	settingsTab    *SettingsTab
	accountInfoTab *AccountsInfoTab
	faxReportTab   *FaxReportTab
	sendFaxTab     *SendFaxTab

	parent *fyne.Window
}

// TabManagementSignal defines a function signature for handling tab management signals.
type TabManagementSignal func(...int)

// NewTabeManagement creates a new instance of TabManagement.
//
// Steps:
// 1. Create a new TabContainer.
// 2. Initialize and return a new TabManagement instance.
//
// Parameters:
//   - apiInst: An instance of the API user interface (ApiUI).
//   - parent: The parent Fyne window.
//
// Returns:
//   - *TabManagement: A pointer to the created TabManagement instance.
func NewTabeManagement(apiInst *api.IApiUICalls, parent *fyne.Window) *TabManagement {
	tabs := widget.NewTabContainer()
	tabManagement := &TabManagement{
		tabs:   tabs,
		parent: parent,
	}
	tabManagement.InitTabs(apiInst)
	return tabManagement
}

// InitTabs initializes the tabs in the TabManagement.
//
// Steps:
// 1. Create instances of SettingsTab, AccountsInfoTab, and FaxReportTab.
// 2. Load UI, data, and set signal functions for each tab.
// 3. Add tabs to the TabContainer based on data availability.
//
// Parameters:
//   - apiInst: An instance of the API user interface (ApiUI).
//
// Returns:
//
//	None
func (m *TabManagement) InitTabs(apiInst *api.IApiUICalls) {
	m.settingsTab = NewSettingsTab(apiInst, m.parent)
	m.settingsTab.initUI()
	m.settingsTab.loadData()
	m.settingsTab.setSignalFunc(m.signalFunc)

	m.accountInfoTab = NewAccountInfoTab(apiInst, m.parent)
	m.accountInfoTab.initUI()

	m.faxReportTab = NewFaxReportTab(apiInst, m.parent)
	m.faxReportTab.initUI()

	m.sendFaxTab = NewSendFaxTab(apiInst, m.parent)
	m.sendFaxTab.initUI()
	m.sendFaxTab.setSignalFunc(m.signalFunc)

	m.tabs = widget.NewTabContainer()

	m.tabs.Items = append(m.tabs.Items, m.settingsTab.GetTab())
	m.tabs.SelectTab(m.settingsTab.GetTab())

	if m.accountInfoTab.IsDataLoaded() {
		m.tabs.Items = append(m.tabs.Items, m.accountInfoTab.GetTab())
		m.tabs.SelectTab(m.accountInfoTab.GetTab())

		if m.faxReportTab.IsDataLoaded() {
			m.tabs.Items = append(m.tabs.Items, m.faxReportTab.GetTab())
		}

		m.tabs.Items = append(m.tabs.Items, m.sendFaxTab.GetTab())
	}
}

// GetTabContainer returns the TabContainer instance.
//
// Parameters:
//
//	None
//
// Returns:
//   - *widget.TabContainer: The TabContainer instance.
func (m *TabManagement) GetTabContainer() *widget.TabContainer {
	return m.tabs
}

// signalFunc handles signals emitted by tabs.
//
// Steps:
// 1. Check the signal ID to determine the type of signal.
// 2. Perform actions based on the signal type, such as updating the displayed tabs.
//
// Parameters:
//   - signals: Variable number of integers representing signal IDs.
//
// Returns:
//
//	None
func (m *TabManagement) signalFunc(signals ...int) {

	if len(signals) == 0 {
		logger.Inst().Info("No signal id available")
		return
	}

	switch signals[0] {
	case SIGNAL_SETTINGS_SAVED:
		m.removeAllTabs()
		if !m.accountInfoTab.IsDataLoaded() {
			forms.ShowError("can not load data from settings data", m.parent)
			m.tabs.SelectTab(m.settingsTab.GetTab())
		} else {
			m.tabs.Items = append(m.tabs.Items, m.accountInfoTab.GetTab())
			m.tabs.Items = append(m.tabs.Items, m.faxReportTab.GetTab())
			m.tabs.Items = append(m.tabs.Items, m.sendFaxTab.GetTab())
			m.tabs.SelectTab(m.accountInfoTab.GetTab())
		}
		break

	case SIGNAL_LOGOUT:
		m.removeAllTabs()
		if !m.accountInfoTab.IsDataLoaded() {
			m.tabs.SelectTab(m.settingsTab.GetTab())
		}
		break

	case SIGNAL_FAX_SENT_SUCCESS:
		m.tabs.SelectTab(m.faxReportTab.GetTab())
		m.faxReportTab.onLoadClick()

	default:
		forms.ShowError(fmt.Sprintf("the signal id %d is not available", signals[0]), m.parent)
		break
	}

}

func (m *TabManagement) removeAllTabs() {
	m.tabs.Remove(m.accountInfoTab.GetTab())
	m.tabs.Remove(m.faxReportTab.GetTab())
	m.tabs.Remove(m.sendFaxTab.GetTab())
}
