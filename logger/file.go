package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	fileLogger := lumberjack.Logger{
		Filename:   "logs/ngecli.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     7,
		Compress:   true,
	}

	writer := zapcore.AddSync(&fileLogger)

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encodeConfig),
		writer,
		zap.DebugLevel,
	)

	core = zapcore.NewTee(
		fileCore,
	)

	logger = zap.New(core)
}
