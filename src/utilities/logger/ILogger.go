package logger

// ILogger is an interface for logging messages.
type ILogger interface {
	// Info logs an informational message.
	// Parameters:
	//   - message: The message to be logged.
	Info(message string)
	// Error logs an error message.
	// Parameters:
	//   - message: The error message to be logged.
	Error(message string)
}
