package log

import "testing"

// Test zap log init function.
func TestInitLogger(t *testing.T) {
	logger := InitLogger("test.log", "debug")
	if logger == nil {
		t.Error("Test InitLogger function failed!")
		return
	}
	t.Log("Test InitLogger function success!")
	logger.Debug("test debug log")
	logger.Info("test info log")
	logger.Error("test err log")
}
