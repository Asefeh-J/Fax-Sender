package impl

import (
	"faxsender/src/utilities"
	"fmt"
	"os"
	"path"
)

const (
	CGO_FLAG_LINUX   = "-ldflags=-linkmode external %v -X main.version=%v"
	CGO_FLAG_WINDOWS = "-ldflags= -H=windowsgui -X %v main.version=%v"
	CGO_STATIC_FLAGS = "-extldflags  -static"

	LINUX_EXTENSION   = ".o"
	WINDOWS_EXTENSION = ".exe"

	MINGW_64BIT_CIMPILER = "x86_64-w64-mingw32-gcc"
	MINGW_32BIT_CIMPILER = "i686-w64-mingw32-gcc"
)

// deployLinux builds a Linux executable.
// Parameters:
//   - fileName: The name of the executable file.
//   - sourcePath: The source code path.
//   - additionalFlags: Additional flags for the build process.
//   - envArgs: Environment variables as variadic parameters.
func deployLinux(fileName string, sourcePath string, additionalFlags string, envArgs ...string) {
	version, err := getVersion()
	if err != nil {
		print(err.Error())
		os.Exit(utilities.ERROR_CODE_VERSION_FILE_NOT_FOUND)
	}
	binPath, _ := utilities.GetExecutablePath()

	linuxExecutablePath := path.Join(binPath, fileName+LINUX_EXTENSION)

	println("software version is: " + version)

	terminalArgs := &utilities.TerminalArgs{
		Command: "go",
		Args:    []string{"build", "-tags", "-mod=vendor", "-o", linuxExecutablePath, fmt.Sprintf(CGO_FLAG_LINUX, additionalFlags, version), sourcePath},
		EnvArgs: envArgs,
	}

	utilities.ExecuteOnTerminalArgs(terminalArgs)
}

// deployWindows builds a Windows executable.
// Parameters:
//   - fileName: The name of the executable file.
//   - sourcePath: The source code path.
func deployWindows(fileName string, sourcePath string, args ...string) {
	version, err := getVersion()
	if err != nil {
		print(err.Error())
		os.Exit(utilities.ERROR_CODE_VERSION_FILE_NOT_FOUND)
	}

	binPath, _ := utilities.GetExecutablePath()

	windowsExceutablePath := ""

	println("software version is: " + version)

	osArchitecture := args[0]
	compiler := MINGW_64BIT_CIMPILER
	switch osArchitecture {
	case "amd64":
		compiler = MINGW_64BIT_CIMPILER
		windowsExceutablePath = path.Join(binPath, fileName+"_64"+WINDOWS_EXTENSION)
	case "386":
		compiler = MINGW_32BIT_CIMPILER
		windowsExceutablePath = path.Join(binPath, fileName+"_32"+WINDOWS_EXTENSION)
	default:
		panic("the compiler for windows can not be recognized.")
	}

	terminalArgs := &utilities.TerminalArgs{
		Command: "go",
		Args:    []string{"build", "-tags", "versiontag", "-mod=vendor", "-o", windowsExceutablePath, fmt.Sprintf(CGO_FLAG_WINDOWS, "", version), sourcePath},
		EnvArgs: []string{
			"GOOS=windows",
			fmt.Sprintf("GOARCH=%s", osArchitecture),
			"CGO_ENABLED=1",
			fmt.Sprintf("CC=%s", compiler),
		},
	}
	utilities.ExecuteOnTerminalArgs(terminalArgs)
}

// DeployLinuxDaemon deploys the Linux daemon executable.
func DeployLinuxDaemon() {
	deployLinux(utilities.APP_EXEC_FILE_NAME, getSourceBackendPath(), CGO_STATIC_FLAGS)
}

// DeployLinuxUi deploys the Linux UI executable.
// Parameters:
//   - architecture: The target architecture.
func DeployLinuxUi(architecture string) {
	deployLinux(utilities.APP_EXEC_UI_FILE_NAME,
		getSourceUIPath(),
		"",
		"GOOS=linux",
		fmt.Sprintf("GOARCH=%s", architecture),
		"CGO_ENABLED=1",
	)
}

// DeployWindowsDaemon deploys the Windows daemon executable.
func DeployWindowsDaemon() {
	deployWindows(utilities.APP_EXEC_UI_FILE_NAME, getSourceBackendPath())
}

// DeployWindowsUi deploys the Windows UI executable.
// Parmeters:
// - arch: the target architecture.
func DeployWindowsUi(arch string) {
	deployWindows(utilities.APP_EXEC_UI_FILE_NAME, getSourceUIPath(), arch)
}

// DeployDebPackages deploys Debian packages.
// Parameters:
//   - arch: The target architecture.
func DeployDebPackages(arch string) {
	version, err := getVersion()
	if err != nil {
		println("error in read version file : " + err.Error())
		os.Exit(utilities.ERROR_CODE_VERSION_FILE_NOT_FOUND)
	}

	fmt.Printf("the version is : %v", version)

	DeployLinuxUi(arch)

	deb := &DebDeployment{}
	dirName := deb.debGenerateDir(version, arch)
	dirAbsPath := deb.createAndGenerateBasePath(dirName)
	deb.moveControlPart(dirAbsPath, arch)
	deb.moveInstFile(dirAbsPath)
	deb.createInstallationPath(dirAbsPath)
	deb.moveExecutableFiles(dirAbsPath)
	deb.copyConfigFileToConfigPath(dirAbsPath)
	deb.createPackage(dirAbsPath)
	deb.removeUnsedFiles(dirAbsPath)
}

// DeployRpmPackage deploys RPM packages.
// Parameters:
//   - arch: The target architecture.
func DeployRpmPackage(arch string) {
	version, err := getVersion()
	panicIfHasError(err)

	fmt.Printf("the version is : %s", version)

	rpm := &RpmDeployment{
		arch:    arch,
		version: version,
	}
	rpm.checkExecutableExist()
	rpm.copySpecFile()
	rpm.copyPostInstFile()
	rpm.initRpmPaths()
	rpm.copyfileToRpmBuild()
	rpm.createRpmPackage()
	rpm.moveBuildFileToBin()

	rpm.removeAllTmpFiles()
}
