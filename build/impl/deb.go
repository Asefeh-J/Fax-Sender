package impl

import (
	"faxsender/src/utilities"
	"fmt"
	"os"
	"path"
	"strings"
)

const (
	APP_NAME_CONTENT     = "#APPNAME#"
	ARCHITECTURE_CONTENT = "#ARCHITECTURE#"
	POSTINST_FILE_NAME   = "postinst"
	POSTRM_FILE_NAME     = "postrm"
	DEBIAN_FILE_NAME     = "DEBIAN"
)

type DebDeployment struct {
}

// debGenerateDir generates the directory structure for the Debian package.
// Parameters:
//   - version: The version of the application.
//   - architecture: The target architecture.
//
// Returns:
//   - The generated directory structure as a string.
func (d *DebDeployment) debGenerateDir(version string, architecture string) string {
	return fmt.Sprintf("%s_%s-1_%s", strings.ToLower(utilities.APP_NAME), strings.Trim(version, "\n"), architecture)
}

// debGenearteConfigPath generates the path for the configuration directory within the Debian package.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
//
// Returns:
//   - The path to the configuration directory as a string.
func (d *DebDeployment) debGenearteConfigPath(dirAbsPath string) string {
	return path.Join(dirAbsPath, "etc", strings.ToLower(utilities.APP_NAME))
}

// debControlFileContent generates the content for the Debian control file.
// Parameters:
//   - arch: The target architecture.
//
// Returns:
//   - The content of the control file as a string.
func (d *DebDeployment) debControlFileContent(arch string) string {
	controlFilePath := path.Join(utilities.GetBuildResequencerPath(), "control")
	fileBytes, err := os.ReadFile(controlFilePath)
	panicIfHasError(err)

	fileStrings := string(fileBytes)
	fileStrings = strings.ReplaceAll(fileStrings, APP_NAME_CONTENT, utilities.APP_NAME)
	fileStrings = strings.ReplaceAll(fileStrings, ARCHITECTURE_CONTENT, arch)
	return fileStrings
}

// moveControlPart moves the control part to the specified directory.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
//   - arch: The target architecture.
func (d *DebDeployment) moveControlPart(dirAbsPath string, arch string) {
	DEBIAN_path := path.Join(dirAbsPath, DEBIAN_FILE_NAME)
	utilities.CreateDirectory(DEBIAN_path)

	controlFileContent := d.debControlFileContent(arch)
	err := utilities.WriteInFile(path.Join(DEBIAN_path, "control"), controlFileContent)
	panicIfHasError(err)
}

// createAndGenerateBasePath creates and generates the base directory for the Debian package.
// Parameters:
//   - dirName: The name of the directory.
//
// Returns:
//   - The absolute path to the main directory as a string.
func (d *DebDeployment) createAndGenerateBasePath(dirName string) string {
	executablePath, err := utilities.GetExecutablePath()
	panicIfHasError(err)

	println("the dir name is : ", dirName)
	dirAbsPath := path.Join(executablePath, dirName)
	utilities.RemoveFolderIfExist(dirAbsPath)
	utilities.CreateDirectory(dirAbsPath)
	return dirAbsPath
}

// createInstallationPath creates necessary directories in the main directory.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) createInstallationPath(dirAbsPath string) {
	sharePath := path.Join(dirAbsPath, "bin")
	configDir := d.debGenearteConfigPath(dirAbsPath)

	utilities.CreateDirectory(sharePath)
	utilities.CreateDirectory(path.Join(configDir, "bin", "logs"))
}

// moveExecutableFiles moves executable files to the specified directory.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) moveExecutableFiles(dirAbsPath string) {
	execDir, err := utilities.GetExecutablePath()
	sharePath := path.Join(dirAbsPath, "bin")
	panicIfHasError(err)

	ui_file_name := utilities.APP_EXEC_UI_FILE_NAME + LINUX_EXTENSION
	ui_file_path := path.Join(execDir, ui_file_name)

	utilities.MoveFile(ui_file_path, path.Join(sharePath, ui_file_name))
}

// moveInstFile moves installation files to the specified directory.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) moveInstFile(dirAbsPath string) {
	debianAbsPath := path.Join(dirAbsPath, DEBIAN_FILE_NAME)
	utilities.CreateDirectory(debianAbsPath)

	postInstAbsPath := path.Join(utilities.GetBuildResequencerPath(), POSTINST_FILE_NAME)
	postinstContent, err := os.ReadFile(postInstAbsPath)
	panicIfHasError(err)

	destinationPostinstAbsPath := path.Join(debianAbsPath, POSTINST_FILE_NAME)

	err = utilities.WriteInFileWithPerm(destinationPostinstAbsPath, string(postinstContent), 0775)
	panicIfHasError(err)

	postrmInstAbsPath := path.Join(utilities.GetBuildResequencerPath(), POSTRM_FILE_NAME)
	postrmContent, err := os.ReadFile(postrmInstAbsPath)
	panicIfHasError(err)

	destinationPostrmPath := path.Join(debianAbsPath, POSTRM_FILE_NAME)

	err = utilities.WriteInFileWithPerm(destinationPostrmPath, string(postrmContent), 0755)
	panicIfHasError(err)
}

// copyConfigFileToConfigPath copies the configuration file to the configuration directory.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) copyConfigFileToConfigPath(dirAbsPath string) {
	config_file_path := path.Join(utilities.GetBuildResequencerPath(), utilities.CONFIG_FILE_NAME)
	config_file_content, err := os.ReadFile(config_file_path)
	panicIfHasError(err)

	base_config_path := d.debGenearteConfigPath(dirAbsPath)
	config_dir := path.Join(base_config_path, "bin", utilities.CONFIG_FILE_NAME)
	utilities.WriteInFile(config_dir, string(config_file_content))
}

// createPackage creates the Debian package.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) createPackage(dirAbsPath string) {
	utilities.ExecuteOnTerminal("dpkg-deb", "--build", "--root-owner-group", dirAbsPath)
}

// removeUnsedFiles removes unused files and directories.
// Parameters:
//   - dirAbsPath: The absolute path to the main directory.
func (d *DebDeployment) removeUnsedFiles(dirAbsPath string) {
	err := utilities.RemoveFolderIfExist(dirAbsPath)
	panicIfHasError(err)
}
