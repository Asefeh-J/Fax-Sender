package tabs

import (
	"faxsender/src/api"
	"faxsender/src/ui/forms"
	"faxsender/src/utilities/logger"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// SettingsTab represents the settings tab in the UI.
type SettingsTab struct {
	hostnameEntry *widget.Entry
	usernameEntry *widget.Entry
	passwordEntry *widget.Entry
	submitButton  *widget.Button
	logoutButton  *widget.Button

	mainContainer *fyne.Container
	tabItem       *container.TabItem

	api    *api.IApiUICalls
	parent *fyne.Window

	signalFunc TabManagementSignal

	ITab
}

// NewSettingsTab creates a new instance of SettingsTab.
//
// Parameters:
//   - apiInst: An instance of the API user interface (ApiUI).
//   - parent: The parent Fyne window.
//
// Returns:
//   - *SettingsTab: A pointer to the created SettingsTab instance.
func NewSettingsTab(apiInst *api.IApiUICalls, parent *fyne.Window) *SettingsTab {
	return &SettingsTab{
		api:    apiInst,
		parent: parent,
	}
}

// initUI initializes the UI components of the SettingsTab.
//
// Steps:
// 1. Create Entry and Button widgets for hostname, username, password, submit, and logout.
// 2. Set icons for submit and logout buttons.
// 3. Create VBox and HBox containers to organize the UI components.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (s *SettingsTab) initUI() {
	s.hostnameEntry = widget.NewEntry()
	s.usernameEntry = widget.NewEntry()
	s.passwordEntry = widget.NewPasswordEntry()
	s.submitButton = widget.NewButton("Login", s.onSubmitClick)
	s.submitButton.Icon = theme.ConfirmIcon()
	s.logoutButton = widget.NewButton("Logout", s.onLogoutClick)
	s.logoutButton.Icon = theme.CancelIcon()

	s.mainContainer = container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Hostname", s.hostnameEntry),
			widget.NewFormItem("Username", s.usernameEntry),
			widget.NewFormItem("Password", s.passwordEntry),
		),
		container.NewHBox(
			s.submitButton,
			s.logoutButton,
		),
	)
}

// loadData loads saved settings from the API and populates the UI components.
//
// Steps:
// 1. Call the API to load saved settings.
// 2. Display an error message if there is an error loading settings.
// 3. Populate the UI components with the loaded settings.
//
// Returns:
//   - bool: True if settings were loaded successfully, false otherwise.
func (s *SettingsTab) loadData() bool {
	savedSettings, err := (*s.api).LoadSettings()
	if err != nil {
		logger.Inst().Error(err.Error())
		forms.ShowError("can't load settings from server", s.parent)
	}

	if savedSettings == nil {
		forms.ShowError("the fetched settings is empty", s.parent)
		return false
	} else {
		s.hostnameEntry.SetText(savedSettings.Hostname)
		s.usernameEntry.SetText(savedSettings.Username)
		s.passwordEntry.SetText(savedSettings.Password)
	}
	return true
}

// setSignalFunc sets the signal function for tab management.
//
// Parameters:
//   - signalFunc: The signal function to set.
//
// Returns:
//
//	None
func (s *SettingsTab) setSignalFunc(signalFunc TabManagementSignal) {
	s.signalFunc = signalFunc
}

// GetTab returns the TabItem for this settings tab.
//
// Parameters:
//
//	None
//
// Returns:
//   - *container.TabItem: The TabItem for this settings tab.
func (s *SettingsTab) GetTab() *container.TabItem {
	if s.tabItem == nil {
		s.tabItem = widget.NewTabItem("Settings", s.mainContainer)
		s.tabItem.Icon = theme.SettingsIcon()
	}
	return s.tabItem
}

// onSubmitClick is the callback function for the submit button.
//
// Steps:
// 1. Get user input from the UI components.
// 2. Call the API to save the settings.
// 3. Display an error message if there is an error saving settings.
// 4. Trigger the SIGNAL_SETTINGS_SAVED signal.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (s *SettingsTab) onSubmitClick() {
	_, err := url.ParseRequestURI(s.hostnameEntry.Text)
	if err != nil {
		logger.Inst().Error(err.Error())
		forms.ShowError("error in the parsed url, the url pattern should be 'http(s)://<your ip address>'", s.parent)
		return
	}

	userData := api.UserData{
		Username: s.usernameEntry.Text,
		Password: s.passwordEntry.Text,
		Hostname: s.hostnameEntry.Text,
	}

	err = (*s.api).SaveSettings(userData)
	if err != nil {
		logger.Inst().Error(err.Error())
		forms.ShowError("error in save settings!", s.parent)
		return
	}

	s.signalFunc(SIGNAL_SETTINGS_SAVED)
}

// onLogoutClick is the callback function for the logout button.
//
// Steps:
// 1. Call the API to perform logout.
// 2. Display an error message if there is an error during the logout process.
// 3. Trigger the SIGNAL_LOGOUT signal.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (s *SettingsTab) onLogoutClick() {
	err := (*s.api).Logout()
	if err != nil {
		logger.Inst().Error(err.Error())
		forms.ShowError("error in logout!", s.parent)
		return
	}
	s.signalFunc(SIGNAL_LOGOUT)
}
