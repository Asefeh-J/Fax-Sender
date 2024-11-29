package forms

import (
	"bytes"
	"encoding/base64"
	"faxsender/src/ui/resources"
	"image"
	"image/png"
	"os"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
)

// InitApp initializes the Fyne app settings, setting the theme to LightTheme.
//
// Parameters:
//   - app: A pointer to the Fyne app to be initialized.
//
// Returns:
//
//	None
func InitApp(app *fyne.App) {
	appTheme := theme.DarkTheme()
	(*app).Settings().SetTheme(appTheme)
}

// InitForms initializes the Fyne window settings.
//
// Steps:
// 1. Set the window to a fixed size.
// 2. Set the close intercept to display a confirmation dialog when closing.
//
// Parameters:
//   - window: A pointer to the Fyne window to be initialized.
//
// Returns:
//
//	None
func InitForms(window *fyne.Window) {
	(*window).SetFixedSize(true)

	(*window).SetCloseIntercept(func() {
		ShowConfirm("Close!", "Are you sure you want to close?", window, func(b bool) {
			if b {
				(*window).Close()
				os.Exit(0)
			}
		})
	})

	iconReader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(resources.PRINTER_ICON_BASE64_ENCODED))
	pngIconImage, _, err := image.Decode(iconReader)
	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = png.Encode(&buff, pngIconImage)
	if err != nil {
		panic(err)
	}
	defer buff.Reset()

	iconResource := fyne.NewStaticResource("printer.png", buff.Bytes())
	(*window).SetIcon(iconResource)
}

// SetSize resizes the Fyne window to the specified width and height.
//
// Parameters:
//   - window: A pointer to the Fyne window to be resized.
//   - width: The width to set for the window.
//   - height: The height to set for the window.
//
// Returns:
//
//	None
func SetSize(window *fyne.Window, width int, height int) {
	(*window).Resize(fyne.NewSize(width, height))
}
