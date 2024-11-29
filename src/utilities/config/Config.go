package config

import (
	"faxsender/src/utilities"
	"os"

	"gopkg.in/yaml.v2"
)

// IConfig is an interface representing the application configuration.

var (
	instConfig IConfig = nil
)

// Config represents the application configuration.
type Config struct {
	PortNumber int  `yaml:"port"`
	Verbose    bool `yaml:"verbose"`
	IConfig    `yaml:"-"`
}

// newConfig creates a new configuration and returns it as an IConfig instance.
// Steps:
// 1. Check if the configuration file exists; if not, create it with default values.
// 2. Read the configuration from the file.
//
// Returns:
//   - IConfig: The application configuration as an IConfig instance.
func newConfig() *IConfig {
	var iconfig IConfig

	err := createConfigIfNotExists()
	if err != nil {
		panic(err)
	}

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	iconfig = config
	return &iconfig
}

// createConfigIfNotExists creates a new configuration file if it doesn't exist.
// Steps:
// 1. Get the path for the system configuration file.
// 2. Check if the file exists; if not, create it with default values.
//
// Returns:
//   - error: An error, if any, encountered during file creation.
func createConfigIfNotExists() error {
	path, err := utilities.GetSystemConfigPath()
	if err != nil {
		return err
	}

	if !utilities.CheckIfFileExists(path) {
		config := &Config{
			PortNumber: utilities.DEFAULT_LISTEN_PORT,
			Verbose:    false,
		}

		bytes, err := yaml.Marshal(config)
		if err != nil {
			return err
		}

		err = os.WriteFile(path, []byte(bytes), 0664)
		if err != nil {
			return err
		}
	}

	return nil
}

// readConfig reads the application configuration from the configuration file.
// Steps:
// 1. Get the path for the system configuration file.
// 2. Read the contents of the file.
// 3. Unmarshal the file contents into a Config struct.
//
// Returns:
//   - *Config: The application configuration as a Config instance.
//   - error: An error, if any, encountered during file reading or unmarshalling.
func readConfig() (*Config, error) {
	var retConfig Config

	path, err := utilities.GetSystemConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &retConfig)
	if err != nil {
		return nil, err
	}

	return &retConfig, nil
}

// Inst returns a global instance of the application configuration.
// If the instance does not exist, it creates a new one.
//
// Returns:
//   - IConfig: The application configuration instance.
func Inst() *IConfig {
	if instConfig == nil {
		instConfig = *newConfig()
	}
	return &instConfig
}

// GetPortNumber returns the port number from the configuration.
//
// Returns:
//   - int: The port number.
func (c Config) GetPortNumber() int {
	return c.PortNumber
}

// GetVerbose returns whether the application is in verbose mode from the configuration.
//
// Returns:
//   - bool: True if the application is in verbose mode, false otherwise.
func (c Config) GetVerbose() bool {
	return c.Verbose
}
