package forms

import (
	"errors"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
)

// ShowInfo displays an information dialog with the provided title and message on the specified Fyne window.
//
// Steps:
// 1. Create a new information dialog using dialog.NewInformation.
// 2. Resize the information dialog to the desired size using infoDialog.Resize.
// 3. Show the information dialog using infoDialog.Show.
//
// Parameters:
//   - title: title of the information dialog.
//   - msg: message to be displayed in the information dialog.
//   - w: pointer to the Fyne Window on which the dialog should be displayed.
func ShowInfo(title string, msg string, w *fyne.Window) {
	infoDialog := dialog.NewInformation(title, msg, *w)
	infoDialog.Resize(fyne.NewSize(400, 100))
	infoDialog.Show()
}

// ShowError displays an error dialog with the provided message on the specified Fyne window.
//
// Steps:
// 1. Create a new error dialog using dialog.NewError.
// 2. Resize the error dialog to the desired size using errorDialog.Resize.
// 3. Show the error dialog using errorDialog.Show.
//
// Parameters:
//   - msg: error message to be displayed in the error dialog.
//   - w: pointer to the Fyne Window on which the dialog
func ShowError(msg string, w *fyne.Window) {
	errorDialog := dialog.NewError(errors.New(msg), *w)
	errorDialog.Resize(fyne.NewSize(400, 100))
	errorDialog.Show()
}

// ShowConfirm displays a confirmation dialog with the specified title and message on the provided Fyne window.
// It also takes a callback function that will be executed based on the user's response.
//
// Steps:
// 1. Create a new confirmation dialog using dialog.NewConfirm.
// 2. Resize the confirmation dialog to the desired size using dialog.Resize.
// 3. Show the confirmation dialog using dialog.Show.
//
// Parameters:
//   - title: title of the confirmation dialog.
//   - msg: message to be displayed in the confirmation dialog.
//   - w: pointer to the Fyne Window on which the dialog should be displayed.
//   - callback: function to be executed based on the user's response (true for OK, false for Cancel).
func ShowConfirm(title string, msg string, w *fyne.Window, callback func(bool)) {
	dialog := dialog.NewConfirm(title, msg, callback, *w)
	dialog.Resize(fyne.NewSize(400, 100))
	dialog.Show()
}
