package main

import (
	"faxsender/src/ui/forms/mainform"
	"faxsender/src/ui/forms/sendfaxform"
	"faxsender/src/utilities"
	"faxsender/src/utilities/logger"
	"flag"
	"os"
)

// main is the entry point of the UI, coordinating the initialization
// of the application through command-line flags and subsequently starting
// either the Fax Sender Form or the Main Form.
//
// Steps:
// 1. Declare command-line flags for showFaxSender, filePath, and workingDir.
// 2. Define and parse command-line flags.
// 3. Change the working directory if a custom working directory is specified.
// 4. Initialize project files using utilities.InitProjectFiles.
// 5. Initialize the application-wide logger using logger.InitLog.
// 6. If showFaxSender flag is set, create and show the Fax Sender Form.
// 7. If showFaxSender flag is not set, create and show the Main Form.
//
// Parameters:
//   - None
//
// Returns:
//
//	This function does not return any values. If an error occurs during the
//	initialization or showing of forms, an error message is printed, and the
//	application exits with an appropriate error code.
func main() {
	var showFaxSender bool
	var filePath string
	var workingDir string

	flag.BoolVar(&showFaxSender, "show-sender", false, "show the fax sender mode")
	flag.StringVar(&filePath, "file-path", "", "the file path to send the fax, this is using in the show-sender mode.")
	flag.StringVar(&workingDir, "working-dir", "", "the working directory to save config.yaml and settings.bin and other settings files.")
	flag.Parse()

	if workingDir != "" {
		os.Chdir(workingDir)
	}

	err := utilities.InitProjectFiles()
	if err != nil {
		println(err)
		os.Exit(utilities.ERROR_CODE_IN_INIT_FILE)
	}

	logger.InitLog()

	if showFaxSender {
		faxSenderForm := sendfaxform.NewSendFaxForm(filePath, nil)
		faxSenderForm.InitControls()
		faxSenderForm.CheckPathWithPanic()
		faxSenderForm.Show()
	} else {
		mainForms := mainform.NewMainForm()
		mainForms.InitControls()
		mainForms.Show()
	}
}
