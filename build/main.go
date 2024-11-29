package main

import (
	"faxsender/build/impl"
	"faxsender/src/utilities"
	"os"
	"path"
)

// main is the entry point of the application and handles various command-line actions.
// Steps:
// 1. Check if there are command-line arguments provided.
// 2. Retrieve the command-line arguments, excluding the program name.
// 3. Perform different actions based on the first argument (command).
//
// Parameters:
//   - None
//
// Returns:
//   - None
func main() {
	if len(os.Args) < 2 {
		println("enter a valid argument to execute...")
		os.Exit(utilities.ERROR_CODE_NOT_ENOUGH_ARGUMENT)
	}

	args := os.Args[1:]

	println(utilities.SEPARATOR)
	switch args[0] {
	case "run_ui_sender":
		println("going to run ui sender mode")
		docsPath, _ := utilities.GetDocsFilePath()
		utilities.ExecuteOnTerminal("go", "run", "./src/ui/main.go", "-show-sender", "-file-path", path.Join(docsPath, "test.docx"))
		break
	case "run_ui":
		println("going to run ui:")
		utilities.ExecuteOnTerminal("go", "run", "./src/ui/main.go")
		break
	case "run":
		println("going to run:")
		utilities.ExecuteOnTerminal("go", "run", "./src/server/main.go")
		break
	case "init":
		println("going to init")
		utilities.ExecuteOnTerminal("mkdir", "-p", "./bin/logs")
		break
	case "package":
		println("going to mod tidy packages...")
		utilities.ExecuteOnTerminal("go", "mod", "tidy")
		println("going to get vendor")
		utilities.ExecuteOnTerminal("go", "mod", "vendor")
		break
	case "deploy_linux_daemon":
		println("going to deploy for linux daemon")
		impl.DeployLinuxDaemon()
		break
	case "deploy_linux_ui":
		println("going to deploy for linux ui")
		arch := args[1]
		impl.DeployLinuxUi(arch)
		break
	case "deploy_windows_daemon":
		println("going to deploy for windows daemon")
		impl.DeployWindowsDaemon()
		break
	case "deploy_windows_ui":
		println("going to deploy for windows ui")
		arch := args[1]
		impl.DeployWindowsUi(arch)
		break
	case "deploy_debs":
		println("going to deploy deb packages")
		arch := args[1]
		impl.DeployDebPackages(arch)
		break
	case "deploy_rpm":
		println("going to deploy rpm package")
		arch := args[1]
		impl.DeployRpmPackage(arch)
		break
	default:
		println(utilities.SEPARATOR)
		println("unknown command! the command : " + args[0])
		os.Exit(utilities.ERROR_CODE_UNKOWN_COMMAND)
		break
	}

	println("DONE!")
}
