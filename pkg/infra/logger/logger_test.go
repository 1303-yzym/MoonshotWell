package logger

import (
	"testing"
)

func TestAppLog(t *testing.T) {
	Logger.Load().AppLog.Info("hello world")
}
