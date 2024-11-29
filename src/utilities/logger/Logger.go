// Logger is a package for handling logging in the application.
// It uses the Zap logger and supports logging to both the console and a log file.
//
// Steps:
// 1. Initialize the logger with the specified configuration.
// 2. Define custom encoders for time and log levels.
// 3. Configure Zap logger with output paths for both console and log file.
// 4. Create a custom core for Zap that includes console and file output.
//
// Parameters:
//
//	None
//
// Returns:
//
//	None
package logger

import (
	"faxsender/src/utilities"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instLog   *Logger
	zapLogger *zap.Logger
	zapConfig zap.Config
)

// Logger represents a logger instance.
type Logger struct {
	ILogger
}

// encodeTime formats the timestamp for logging.
func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("[%s]", t.Format("2006-01-02 15:04:05")))
}

// encodeLevel formats the log level for logging.
func encodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// InitLog initializes the logger with a specific configuration.
func InitLog() {
	var err error
	cfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          &zap.SamplingConfig{},
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "debug",
			EncodeLevel:   zapcore.CapitalLevelEncoder,
			TimeKey:       "timestamp",
			EncodeTime:    zapcore.EpochMillisTimeEncoder,
			StacktraceKey: "stack",
			LineEnding:    "\n",
		},
		OutputPaths:      []string{"stdout", path.Join(utilities.GetLogsPath(), "logs.log")},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields:    map[string]interface{}{},
	}

	cfg.EncoderConfig.EncodeTime = encodeTime
	cfg.EncoderConfig.EncodeLevel = encodeLevel
	zapConfig = cfg
	zapLogger, err = zapConfig.Build(zap.WrapCore(zapCore))
	if err != nil {
		fmt.Println(err)
		os.Exit(utilities.ERROR_CODE_INIT_LOG_ERROR)
	}
	defer zapLogger.Sync()
}

// zapCore creates a custom Zap core that supports console and log file output.
func zapCore(c zapcore.Core) zapcore.Core {
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   zapConfig.OutputPaths[1],
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     100, //days
		Compress:   true,
	})

	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zapConfig.EncoderConfig),
		w,
		zap.DebugLevel,
	)

	pe := zap.NewProductionEncoderConfig()
	pe.EncodeTime = encodeTime
	pe.EncodeLevel = encodeLevel
	consoleEncoder := zapcore.NewConsoleEncoder(pe)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	cores := zapcore.NewTee(c, fileCore, consoleCore)

	return cores
}

// Inst returns the logger instance.
func Inst() ILogger {
	if instLog == nil {
		instLog = &Logger{}
	}
	return instLog
}

// Info logs an informational message.
func (l *Logger) Info(message string) {
	zapLogger.Info(message)
}

// Error logs an error message.
func (l *Logger) Error(message string) {
	zapLogger.Error(message)
}
