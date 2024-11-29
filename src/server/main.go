package main

import (
	"faxsender/src/api"
	"faxsender/src/utilities"
	"faxsender/src/utilities/config"
	"faxsender/src/utilities/logger"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	workingDir string
)

// Init initializes the application by performing the following steps:
//
// 1. Setting up the working directory based on command-line flags using InitWorkingDir.
// 2. Initializing project-specific files with utilities.InitProjectFiles.
// 3. Configuring application-wide logging with InitLogConfig.
// 4. Initializing system configuration and logging the server port using InitSystemConfig.
//
// This function serves as a centralized entry point for initializing various aspects
// of the application, making it easier to manage and understand the startup process.
func Init() {
	InitWorkingDir()
	utilities.InitProjectFiles()
	InitLogConfig()
	InitSystemConfig()
}

// InitWorkingDir sets up the working directory based on the command-line flags.
//
// This function uses the "flag" package to parse the command-line arguments and
// retrieve the value provided for the "working-dir" flag. If a non-empty directory
// path is provided, it changes the current working directory to that path. If the
// specified directory does not exist, it attempts to create the directory and exits
// the application with an error code if the creation fails.
func InitWorkingDir() {
	flag.StringVar(&workingDir, "working-dir", "", "the directory to work with")
	flag.Parse()

	if workingDir != "" {
		os.Chdir(workingDir)
		err := os.MkdirAll(workingDir, os.ModeAppend)
		if err != nil {
			fmt.Print(err)
			os.Exit(utilities.ERROR_CODE_WORKING_DIR_NOT_FOUND)
		}
	}
}

// InitSystemConfig initializes the system configuration and logs the server port.
//
// This function retrieves the server configuration using the config.Inst() function,
// logs an informational message using the application-wide logger, indicating the
// configured server port. The message includes the actual server port retrieved
// from the configuration.
func InitSystemConfig() {
	cfg := *config.Inst()
	logger.Inst().Info(fmt.Sprintf("the server port is: %v", cfg.GetPortNumber()))
}

// InitLogConfig initializes the logging configuration.
//
// This function is responsible for initializing the application-wide logging.
// It calls the InitLog function from the logger package to set up the logging
// configuration. After successful initialization, it logs an informational
// message using the application-wide logger, indicating that the logger has
// been successfully initialized.
func InitLogConfig() {
	logger.InitLog()
	logger.Inst().Info("logger initialized")
}

// StartServer initializes the Gin router, sets up API routes, and starts the server.
//
// This function retrieves the server configuration, initializes a Gin router,
// sets up API routes using the InitRouters function from the api package,
// and then starts the server by calling the listenWithoutCertificates function.
//
// It logs an informational message indicating the port on which the server is
// about to listen before initiating the server startup process.
func StartServer() {
	cfg := *config.Inst()
	port := cfg.GetPortNumber()

	router := gin.Default()
	api.InitRouters(router)

	logger.Inst().Info(fmt.Sprintf("going to listen on port %d", port))
	listenWithoutCertificates(router, port)
}

// listenWithoutCertificates starts the server to listen on the specified port.
//
// Parameters:
//   - router: A pointer to a Gin Engine instance, which represents the HTTP router.
//   - port: An integer specifying the port on which the server will listen.
//
// Returns:
//
//	This function does not return any values. If an error occurs during server startup,
//	an error message is logged using the application-wide logger.
func listenWithoutCertificates(router *gin.Engine, port int) {
	err := router.Run(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Inst().Error(fmt.Sprintf("failed to start the server:%v", err))
		return
	}
}

// main is the entry point of the application, coordinating the initialization
// of the application through the Init function and subsequently starting
// the server using the StartServer function.
func main() {
	Init()
	StartServer()
}
