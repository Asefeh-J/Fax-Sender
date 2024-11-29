package config

// IConfig is an interface representing the application configuration.
type IConfig interface {
	// GetPortNumber retrieves the port number from the configuration.
	// Returns:
	//   - int: The port number.
	GetPortNumber() int
	// GetVerbose checks whether the application is in verbose mode from the configuration.
	// Returns:
	//   - bool: True if the application is in verbose mode, false otherwise.
	GetVerbose() bool
}
