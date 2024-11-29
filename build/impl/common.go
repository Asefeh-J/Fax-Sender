package impl

import (
	"faxsender/src/utilities"
	"os"
	"path"
)

// getVersion retrieves the application's version from a version.txt file.
// Steps:
// 1. Get the build path using utilities.GetBuildPath().
// 2. Read the version.txt file in the build path.
// 3. Clean up the file path using utilities.CleanupFilePath().
//
// Returns:
//   - The application's version as a string.
//   - An error if any.
func getVersion() (string, error) {
	buildPath, err := utilities.GetBuildPath()
	if err != nil {
		println("Failed to get build path")
		return "", err
	}

	bytes, err := os.ReadFile(path.Join(buildPath, "version.txt"))
	return utilities.CleanupFilePath(string(bytes)), err
}

// getSourceBackendPath returns the file path for the main.go file in the server directory of the source code.
// Returns:
//   - The file path as a string.
func getSourceBackendPath() string {
	srcPath, _ := utilities.GetSourcePath()
	return path.Join(srcPath, "server", "main.go")
}

// getSourceUIPath returns the file path for the main.go file in the UI directory of the source code.
// Returns:
//   - The file path as a string.
func getSourceUIPath() string {
	srcPath, _ := utilities.GetSourcePath()
	srcPath = path.Join(srcPath, "ui")
	return path.Join(srcPath, "main.go")
}

// panicIfHasError panics if the given error is not nil.
// Parameters:
//   - err: The error to check.
func panicIfHasError(err error) {
	if err != nil {
		panic(err)
	}
}
