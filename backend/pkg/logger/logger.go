package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log *zap.Logger

// Init initializes the logger
func Init() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	// Parse the log level
	var level zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure console encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create console core
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		level,
	)

	// Create logger
	log = zap.New(consoleCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

// Debug logs a debug message
func Debug(message string, fields ...zapcore.Field) {
	if log == nil {
		Init()
	}
	log.Debug(message, fields...)
}

// Info logs an info message
func Info(message string, fields ...zapcore.Field) {
	if log == nil {
		Init()
	}
	log.Info(message, fields...)
}

// Warn logs a warning message
func Warn(message string, fields ...zapcore.Field) {
	if log == nil {
		Init()
	}
	log.Warn(message, fields...)
}

// Error logs an error message
func Error(message string, fields ...zapcore.Field) {
	if log == nil {
		Init()
	}
	log.Error(message, fields...)
}

// Fatal logs a fatal message and exits
func Fatal(message string, fields ...zapcore.Field) {
	if log == nil {
		Init()
	}
	log.Fatal(message, fields...)
}

// Field creates a field for structured logging
func Field(key string, value interface{}) zapcore.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case bool:
		return zap.Bool(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.String(key, fmt.Sprintf("%v", v))
	}
}

// Close flushes any buffered log entries
func Close() {
	if log != nil {
		_ = log.Sync()
	}
}
