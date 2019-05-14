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

// Sugar converts a Logger to a SugaredLogger.
func Sugar() *zap.SugaredLogger { return logger.Sugar() }

// Named adds a new path segment to the logger's name. Segments are joined by
// periods.
func Named(name string) *zap.Logger { return logger.Named(name) }

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return logger.WithOptions(opts...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zapcore.Field) *zap.Logger { return logger.With(fields...) }

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

// DPanic log panic level message & throw a panic in development mode
func DPanic(msg string, fields ...zapcore.Field) {
	logger.DPanic(msg, fields...)
}

// Panic log panic level message & throw a panic
func Panic(msg string, fields ...zapcore.Field) {
	logger.Panic(msg, fields...)
}

// Fatal log fatal level message & exit
func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}
