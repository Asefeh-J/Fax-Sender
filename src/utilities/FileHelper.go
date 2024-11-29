package utilities

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	TEXT_FILE_EXTENSION  string = "txt"
	PDF_FILE_EXTENSION   string = "pdf"
	TIFF_FILE_EXTENSION  string = "tiff"
	TIF_FILE_EXTENSION   string = "tif"
	JPEG_FILE_EXTENSION  string = "jpeg"
	JPG_FILE_EXTENSION   string = "jpg"
	PNG_FILE_EXTENSION   string = "png"
	DOC_FILE_EXTENSION   string = "doc"
	DOCX_FILE_EXTENSION  string = "docx"
	ODT_FILE_EXTENSION   string = "odt"
	EMPTY_FILE_EXTENSION string = ""

	MS_WORD_CONTENT_TYPE string = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

// GetExecutablePath returns the path to the directory where the executable is located.
//
// Returns:
//   - string: The path to the executable directory.
//   - error: An error if the path cannot be determined.
func GetExecutablePath() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(p, "bin"), nil
}

// GetBuildPath returns the path to the build directory.
//
// Returns:
//   - string: The path to the build directory.
//   - error: An error if the path cannot be determined.
func GetBuildPath() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(p, "build"), nil
}

// GetBuildResequencerPath returns the path to the build resources directory.
//
// Returns:
//   - string: The path to the build resources directory.
func GetBuildResequencerPath() string {
	p, err := GetBuildPath()
	if err != nil {
		panic(err)
	}
	return path.Join(p, "resources")
}

// GetLogsPath returns the path to the logs directory within the executable directory.
//
// Returns:
//   - string: The path to the logs directory.
func GetLogsPath() string {
	binPath, err := GetExecutablePath()
	if err != nil {
		println(err)
		os.Exit(ERROR_CODE_WORKING_DIR_NOT_FOUND)
	}
	return path.Join(binPath, "logs")
}

// GetSourcePath returns the path to the source code directory.
//
// Returns:
//   - string: The path to the source code directory.
//   - error: An error if the path cannot be determined.
func GetSourcePath() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(p, "src"), nil
}

// GetSystemConfigPath returns the path to the system configuration file.
//
// Returns:
//   - string: The path to the system configuration file.
//   - error: An error if the path cannot be determined.
func GetSystemConfigPath() (string, error) {
	exec, err := GetExecutablePath()

	if err != nil {
		return "", err
	}

	return path.Join(exec, CONFIG_FILE_NAME), nil
}

// GetSystemSettingsPath returns the path to the system settings file.
//
// Returns:
//   - string: The path to the system settings file.
//   - error: An error if the path cannot be determined.
func GetSystemSettingsPath() (string, error) {
	exec, err := GetExecutablePath()

	if err != nil {
		return "", err
	}

	return path.Join(exec, SETTINGS_FILE_NAME), nil
}

// CheckIfFileExists checks if a file exists at the specified path.
//
// Parameters:
//   - path: The path to the file to be checked.
//
// Returns:
//   - bool: true if the file exists, false otherwise.
func CheckIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetDocsFilePath returns the path to the "docs" directory.
//
// Returns:
//   - string: The path to the "docs" directory.
//   - error: An error if the path cannot be determined.
func GetDocsFilePath() (string, error) {
	p, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(p, "docs"), nil
}

// GetContentType returns the content type (MIME type) based on the file extension.
//
// Parameters:
//   - extension: The file extension without a dot (e.g., "pdf").
//
// Returns:
//   - string: The corresponding content type.
func GetContentType(extension string) string {
	contentTypes := map[string]string{
		PDF_FILE_EXTENSION:  "application/pdf",
		JPEG_FILE_EXTENSION: "image/jpeg",
		JPG_FILE_EXTENSION:  "image/jpg",
		PNG_FILE_EXTENSION:  "image/png",
		TIFF_FILE_EXTENSION: "image/tiff",
		TIF_FILE_EXTENSION:  "image/tiff",
		DOCX_FILE_EXTENSION: MS_WORD_CONTENT_TYPE,
		DOC_FILE_EXTENSION:  MS_WORD_CONTENT_TYPE,
		ODT_FILE_EXTENSION:  "application/vnd.oasis.opendocument.text",
	}

	return contentTypes[strings.ToLower(extension)]
}

// CalculateFileExtension calculates the file extension from a file path.
//
// Parameters:
//   - fpath: The file path.
//
// Returns:
//   - string: The calculated file extension (e.g., "pdf").
func CalculateFileExtension(fpath string) string {
	fpath = RemoveDot(filepath.Ext(fpath))
	return CleanupFilePath(fpath)
}

// ExtractFileExtension extracts and returns the file extension from a file path.
// If the extension is not recognized, it returns an empty string.
//
// Parameters:
//   - filePath: The file path.
//
// Returns:
//   - string: The extracted file extension (e.g., "pdf").
func ExtractFileExtension(filePath string) string {
	extension := CalculateFileExtension(filePath)

	if extension == TEXT_FILE_EXTENSION ||
		extension == PDF_FILE_EXTENSION ||
		extension == TIFF_FILE_EXTENSION ||
		extension == JPEG_FILE_EXTENSION ||
		extension == JPG_FILE_EXTENSION ||
		extension == PNG_FILE_EXTENSION ||
		extension == DOCX_FILE_EXTENSION ||
		extension == DOC_FILE_EXTENSION ||
		extension == TIF_FILE_EXTENSION ||
		extension == ODT_FILE_EXTENSION {
		return extension
	} else {
		return EMPTY_FILE_EXTENSION
	}
}

// InitProjectFiles initializes project directories and files.
//
// Returns:
//   - error: An error if initialization fails.
func InitProjectFiles() error {
	basePath := GetLogsPath()

	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// CreateDirectory creates a directory with the specified path.
//
// Parameters:
//   - dir: The path of the directory to be created.
//
// Returns:
//   - error: An error if directory creation fails.
func CreateDirectory(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// RemoveFolderIfExist removes a directory and its contents if it exists.
//
// Parameters:
//   - dirPath: The path of the directory to be removed.
//
// Returns:
//   - error: An error if removal fails or the directory does not exist.
func RemoveFolderIfExist(dirPath string) error {
	if CheckIfFileExists(dirPath) {
		return os.RemoveAll(dirPath)
	}
	return errors.New("the directory path does not exist")
}

// MoveFile renames a source file to a destination file.
//
// Parameters:
//   - sourceFile: The source file path.
//   - destinationFile: The destination file path.
//
// Returns:
//   - error: An error if the renaming operation fails.
func MoveFile(sourceFile, destintionFile string) error {
	// return os.Rename(sourceFile, destintionFile)
	err := CopyFile(sourceFile, destintionFile)
	if err != nil {
		return err
	}

	return os.Remove(sourceFile)
}

// CopyFile copies a source file to a destination file.
//
// Parameters:
//   - sourceFile: The source file path.
//   - destinationFile: The destination file path.
//
// Returns:
//   - error: An error if the copying operation fails.
func CopyFile(sourceFile, destinationFile string) error {
	if !CheckIfFileExists(sourceFile) {
		return fmt.Errorf("the source file '%s' does not exist", sourceFile)
	}

	dest, err := os.OpenFile(destinationFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer dest.Close()

	sourceFileBytes, err := os.ReadFile(sourceFile)
	if err != nil {
		return err
	}

	_, err = dest.Write(sourceFileBytes)
	if err != nil {
		return err
	}

	return nil
}

// WriteInFile writes content to a file with the specified path and permissions.
//
// Parameters:
//   - filePath: The path of the file to be written to.
//   - contents: The content to be written to the file.
//   - permission: The file permission (e.g., 0644).
//
// Returns:
//   - error: An error if the write operation fails.
func WriteInFile(filePath, contents string) error {
	return WriteInFileWithPerm(filePath, contents, 0644)
}

// WriteInFileWithPerm writes content to a file with the specified path and permissions.
//
// Parameters:
//   - filePath: The path of the file to be written to.
//   - contents: The content to be written to the file.
//   - permission: The file permission (e.g., 0644).
//
// Returns:
//   - error: An error if the write operation fails.
func WriteInFileWithPerm(filePath, contents string, permission fs.FileMode) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, permission)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(contents))
	return err
}

// AllValidExtensions returns a list of all valid file extensions.
//
// Returns:
//   - []string: A list of valid file extensions (e.g., [".txt", ".pdf"]).
func AllValidExtensions() []string {
	list := make([]string, 0)
	list = append(list, appendDot(TEXT_FILE_EXTENSION))
	list = append(list, appendDot(PDF_FILE_EXTENSION))
	list = append(list, appendDot(TIFF_FILE_EXTENSION))
	list = append(list, appendDot(JPEG_FILE_EXTENSION))
	list = append(list, appendDot(JPG_FILE_EXTENSION))
	list = append(list, appendDot(PNG_FILE_EXTENSION))
	list = append(list, appendDot(DOC_FILE_EXTENSION))
	list = append(list, appendDot(DOCX_FILE_EXTENSION))
	list = append(list, appendDot(TIF_FILE_EXTENSION))
	list = append(list, appendDot(ODT_FILE_EXTENSION))

	return list
}

// CleanupFilePath removes newline characters and extra whitespaces from a string.
//
// Parameters:
//   - str: The string to be cleaned.
//
// Returns:
//   - string: The cleaned string.
func CleanupFilePath(str string) string {
	str = strings.ReplaceAll(str, "\n", "")
	return strings.TrimSpace(str)
}

// appendDot adds a dot (.) to the beginning of a name to form a file extension.
//
// Parameters:
//   - name: The name (e.g., "txt").
//
// Returns:
//   - string: The file extension (e.g., ".txt").
func appendDot(name string) string {
	return fmt.Sprintf(".%v", name)
}

// RemoveDot removes the dot (.) from a name, typically used for cleaning file extensions.
//
// Parameters:
//   - name: The name with or without a dot (e.g., ".txt" or "txt").
//
// Returns:
//   - string: The name without a dot (e.g., "txt").
func RemoveDot(name string) string {
	return strings.Trim(strings.ReplaceAll(name, ".", ""), " ")
}
