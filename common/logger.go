package common

import (
	"github.com/uber-go/zap"
)

// Logger common logger for
var Logger *zap.Logger

// SetLogger to set local logger by caller
func SetLogger(l *zap.Logger) {
	if l == nil {
		panic("logger cannot be nil.")
	}

	Logger = l
}
