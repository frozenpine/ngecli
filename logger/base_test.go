package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	testLogger := Named("test")

	testLogger.Info("this is info", zap.Int("count", 1))

	testLogger.Error("error occoured", zap.Bool("success", false))
}
