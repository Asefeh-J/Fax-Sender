package forms

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

// SettingsForm is a structure that includes a window field for managing
// the user interface and implements the IBaseForm interface.
type SettingsForm struct {
	window *fyne.Window
	IBaseForm
}

// NewSettingsForm creates a new SettingsForm instance using the specified Fyne app.
//
// Steps:
// 1. Create a new Fyne window with the title "login data".
// 2. Initialize the window settings and size using InitForms and SetSize.
// 3. Create and return a new SettingsForm instance.
//
// Parameters:
//   - app: Fyne app instance used to create the window.
//
// Returns:
//   - a new SettingsForm instance.
func NewSettingsForm(app *fyne.App) *SettingsForm {
	window := (*app).NewWindow("login data")

	InitForms(&window)
	SetSize(&window, 300, 200)
	return &SettingsForm{
		window: &window,
	}
}

// InitControls initializes the UI controls for the SettingsForm.
//
// Steps:
//  1. Create an API key entry field using widget.NewEntry.
//  2. Set the content of the window to a vertical container containing:
//     a. A label asking the user to enter the API key.
//     b. The API key entry field.
//     c. A submit button that, when clicked, shows an info dialog with the entered API key.
func (f *SettingsForm) InitControls() {

	apikeyEntry := widget.NewEntry()

	(*f.window).SetContent(container.NewVBox(
		widget.NewLabel("Enter your API key:"),
		apikeyEntry,
		widget.NewButton("Submit", func() {
			apikey := apikeyEntry.Text
			//TODO: Send the API key to the daemon to save it
			ShowInfo("data", apikey, f.window)
		}),
	))
}

// Show displays the SettingsForm window.
func (f *SettingsForm) Show() {
	(*f.window).Show()
}
