package tabs

import (
	"faxsender/src/api"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

// AccountsInfoTab is a tab for displaying account information.
type AccountsInfoTab struct {
	firstNameLabel *widget.Label
	lastNameLabel  *widget.Label
	emailLabel     *widget.Label
	phoneLabel     *widget.Label
	countryLabel   *widget.Label
	companyLabel   *widget.Label

	mainContainer *fyne.Container
	tabItem       *container.TabItem

	api    *api.IApiUICalls
	parent *fyne.Window

	ITab
}

// NewAccountInfoTab creates a new instance of AccountsInfoTab.
//
// Parameters:
//   - apiInst: An instance of the ApiUI for API interactions.
//   - parent: The parent window associated with the tab.
//
// Returns:
//   - *AccountsInfoTab: The created AccountsInfoTab instance.
func NewAccountInfoTab(apiInst *api.IApiUICalls, parent *fyne.Window) *AccountsInfoTab {
	return &AccountsInfoTab{
		api:    apiInst,
		parent: parent,
	}
}

// createText formats a label and its corresponding value into a string.
//
// Parameters:
//   - label: The label describing the information.
//   - value: The value of the information.
//
// Returns:
//   - string: The formatted string combining the label and value.
func (a *AccountsInfoTab) createText(label string, value string) string {
	return fmt.Sprintf("%s: %s", label, value)
}

// fillLabels updates the UI labels with the account information.
//
// Parameters:
//   - accountInfo: An instance of AccountInfo containing user data.
//
// Returns:
//
//	None
func (s *AccountsInfoTab) fillLabels(accountInfo *api.AccountInfo) {
	s.firstNameLabel.Text = s.createText("Firstname", accountInfo.FirstName)
	s.lastNameLabel.Text = s.createText("Lastname", accountInfo.LastName)
	s.emailLabel.Text = s.createText("Email", accountInfo.Email)
	s.phoneLabel.Text = s.createText("Phone", accountInfo.Phone)
	s.countryLabel.Text = s.createText("Country", accountInfo.Country)
	s.companyLabel.Text = s.createText("Company", accountInfo.Country)
}

// initUI initializes the UI components of the AccountsInfoTab.
//
// Steps:
// 1. Create and configure UI components such as labels.
// 2. Set up the layout of the components using containers.
// 3. Populate the labels with default or empty data.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
func (a *AccountsInfoTab) initUI() {
	a.firstNameLabel = widget.NewLabel("")
	a.lastNameLabel = widget.NewLabel("")
	a.emailLabel = widget.NewLabel("")
	a.phoneLabel = widget.NewLabel("")
	a.countryLabel = widget.NewLabel("")
	a.companyLabel = widget.NewLabel("")

	a.mainContainer = container.NewVBox(
		container.NewGridWithColumns(2,
			a.firstNameLabel,
			a.lastNameLabel,
			a.emailLabel,
			a.phoneLabel,
			a.countryLabel,
			a.companyLabel,
		),
	)

	a.fillLabels(api.GetEmptyAccountInfo())
}

// loadData fetches and loads the account information from the API.
//
// Steps:
// 1. Call the API to retrieve the account information.
// 2. Handle any errors during the data retrieval process.
//
// Returns:
//   - bool: True if data is loaded successfully, false otherwise.
func (a *AccountsInfoTab) loadData() bool {
	account_info, err := (*a.api).GetAccountInfo()
	if err != nil {
		return false
	}

	a.fillLabels(account_info)
	return true
}

// GetTab returns the TabItem associated with the AccountsInfoTab.
//
// Steps:
// 1. Return the TabItem associated with the tab.
//
// Returns:
//   - *container.TabItem: The TabItem associated with the tab.
func (a *AccountsInfoTab) GetTab() *container.TabItem {
	if a.tabItem == nil {
		a.tabItem = container.NewTabItem("Account Information", a.mainContainer)
		a.tabItem.Icon = theme.InfoIcon()
	}
	return a.tabItem
}

// IsDataLoaded checks if the account data is loaded.
//
// Returns:
//   - bool: True if the data is loaded, false otherwise.
func (a *AccountsInfoTab) IsDataLoaded() bool {
	return a.loadData()
}
