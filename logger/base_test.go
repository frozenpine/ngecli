package logger

import (
	"testing"

	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	logger.Info("this is info", zap.Int("count", 1))

	logger.Error("error occoured", zap.Bool("success", false))
}
