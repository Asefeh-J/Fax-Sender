package impl

import (
	"faxsender/src/utilities"
	"fmt"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	EMPTY_FILE_NAME                     = ""
	RPM_SPEC_FILE_NAME                  = "print2fax.spec"
	RPM_SPEC_FILE_CONTENT_NAME          = "##APP_NAME##"
	RPM_SPEC_FILE_CONTENT_EXEC_NAME     = "##EXEC_NAME##"
	RPM_SPEC_FILE_CONTENT_ARCH          = "##ARCH##"
	RPM_SPEC_FILE_CONTENT_VERSION       = "##VERSION##"
	RPM_SPEC_FILE_CONTENT_RELEASE       = "##RELEASE##"
	RPM_SPEC_FILE_CONTENT_POSTINST_FILE = "##POSTINST_FILE##"
	RPM_COMMAND                         = "rpmbuild"
	RPM_NO_ARCH                         = "noarch"
	RPM_RELEASE                         = "1"
)

type RpmDeployment struct {
	specFileFullAddress      string
	executableFilePath       string
	rpmBuildBaseFullAddress  string
	rpmBuildBUILDFullAddress string
	rpmBuildSPECFullAddress  string
	rpmBuildRPMSFullAddress  string
	postInstFullAddress      string
	arch                     string
	version                  string
}

// getBinPathWithPanic returns the executable path or panics if an error occurs.
// Returns:
//   - The executable path as a string.
func (r *RpmDeployment) getBinPathWithPanic() string {
	binPath, err := utilities.GetExecutablePath()
	panicIfHasError(err)

	return binPath
}

// executableFileName returns the name of the executable file.
// Returns:
//   - The executable file name as a string.
func (r *RpmDeployment) executableFileName() string {
	return fmt.Sprintf("%s%s", utilities.APP_EXEC_UI_FILE_NAME, LINUX_EXTENSION)
}

// packageName returns the name of the RPM package.
// Returns:
//   - The RPM package name as a string.
func (r *RpmDeployment) packageName() string {
	return fmt.Sprintf("%s-%s-%s.%s.rpm", strings.ToLower(utilities.APP_NAME), r.version, RPM_RELEASE, RPM_NO_ARCH)
}

// checkExecutableExist checks if the executable file exists and panics if it doesn't.
func (r *RpmDeployment) checkExecutableExist() {
	binPath := r.getBinPathWithPanic()
	r.executableFilePath = path.Join(binPath, r.executableFileName())
	if !utilities.CheckIfFileExists(r.executableFilePath) {
		panic("error : the executable file does not exist")
	}
}

// copySpecFile copies the RPM spec file and replaces placeholders with values.
func (r *RpmDeployment) copySpecFile() {
	resourcePath := utilities.GetBuildResequencerPath()
	specFileInResource := path.Join(resourcePath, RPM_SPEC_FILE_NAME)

	fileBytes, err := os.ReadFile(specFileInResource)
	panicIfHasError(err)

	fileStrings := string(fileBytes)
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_NAME, strings.ToLower(utilities.APP_NAME))
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_EXEC_NAME, r.executableFileName())
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_ARCH, RPM_NO_ARCH)
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_VERSION, r.version)
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_RELEASE, RPM_RELEASE)
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_POSTINST_FILE, POSTINST_FILE_NAME)

	r.specFileFullAddress = path.Join(r.getBinPathWithPanic(), RPM_SPEC_FILE_NAME)

	err = utilities.WriteInFile(r.specFileFullAddress, fileStrings)
	panicIfHasError(err)
}

// copyPostInstFile copies the post-installation script file and replaces placeholders with values.
func (r *RpmDeployment) copyPostInstFile() {
	resourcePath := utilities.GetBuildResequencerPath()
	print("the resource path is : ", resourcePath)
	fileBytes, err := os.ReadFile(path.Join(resourcePath, POSTINST_FILE_NAME))
	panicIfHasError(err)

	fileStrings := string(fileBytes)
	fileStrings = strings.ReplaceAll(fileStrings, RPM_SPEC_FILE_CONTENT_EXEC_NAME, r.executableFileName())

	r.postInstFullAddress = path.Join(r.getBinPathWithPanic(), POSTINST_FILE_NAME)

	err = utilities.WriteInFile(r.postInstFullAddress, fileStrings)
	panicIfHasError(err)
}

// initRpmPaths initializes RPM-related paths.
func (r *RpmDeployment) initRpmPaths() {
	err := utilities.ExecuteOnTerminal("rpmdev-setuptree")
	panicIfHasError(err)

	user, err := user.Current()
	panicIfHasError(err)

	r.rpmBuildBaseFullAddress = path.Join(user.HomeDir, RPM_COMMAND)
	r.rpmBuildBUILDFullAddress = path.Join(r.rpmBuildBaseFullAddress, "BUILD")
	r.rpmBuildSPECFullAddress = path.Join(r.rpmBuildBaseFullAddress, "SPECS")
	r.rpmBuildRPMSFullAddress = path.Join(r.rpmBuildBaseFullAddress, "RPMS")
}

// copyfileToRpmBuild copies necessary files to the RPM build directory.
func (r *RpmDeployment) copyfileToRpmBuild() {
	err := utilities.CopyFile(r.specFileFullAddress, path.Join(r.rpmBuildSPECFullAddress, RPM_SPEC_FILE_NAME))
	panicIfHasError(err)

	err = utilities.CopyFile(r.executableFilePath, path.Join(r.rpmBuildBUILDFullAddress, r.executableFileName()))
	panicIfHasError(err)

	err = utilities.CopyFile(r.postInstFullAddress, path.Join(r.rpmBuildBUILDFullAddress, POSTINST_FILE_NAME))
	panicIfHasError(err)
}

// createRpmPackage creates the RPM package.
func (r *RpmDeployment) createRpmPackage() {
	err := utilities.ExecuteOnTerminal(RPM_COMMAND, "-ba", path.Join(r.rpmBuildSPECFullAddress, RPM_SPEC_FILE_NAME))
	panicIfHasError(err)
}

// moveBuildFileToBin moves the built RPM package to the executable path.
func (r *RpmDeployment) moveBuildFileToBin() {
	sourceFile := path.Join(r.rpmBuildRPMSFullAddress, RPM_NO_ARCH, r.packageName())

	binPath, err := utilities.GetExecutablePath()
	panicIfHasError(err)

	err = utilities.MoveFile(sourceFile, path.Join(binPath, r.packageName()))
	panicIfHasError(err)
}

// removeAllTmpFiles removes temporary files and directories.
func (r *RpmDeployment) removeAllTmpFiles() {
	err := utilities.RemoveFolderIfExist(r.specFileFullAddress)
	panicIfHasError(err)

	err = utilities.RemoveFolderIfExist(r.rpmBuildBaseFullAddress)
	panicIfHasError(err)

	err = utilities.RemoveFolderIfExist(r.postInstFullAddress)
	panicIfHasError(err)
}
