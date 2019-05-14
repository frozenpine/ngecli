package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	highPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	lowPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	})

	encodeConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "lvl",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	core zapcore.Core

	logger *zap.Logger
)

func init() {
	consoleOut := zapcore.Lock(os.Stdout)
	consoleErr := zapcore.Lock(os.Stderr)

	consoleEncoder := zapcore.NewConsoleEncoder(encodeConfig)

	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErr, highPriority),
		zapcore.NewCore(consoleEncoder, consoleOut, lowPriority),
	)

	logger = zap.New(core)
}

// Flush make sure log in buffer will be synced.
func Flush() {
	if logger == nil {
		return
	}

	// logger.D
}

// Debug log debug level message
func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

// Info log info level message
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// Warn log warn level message
func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

// Error log error level message
func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}
