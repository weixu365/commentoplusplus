package util

import (
	"testing"
)

func TestLoggerCreateBasics(t *testing.T) {
	logger = GetLogger()

	if logger == nil {
		t.Errorf("logger null after GetLogger()")
		return
	}

	logger.Debugf("test message please ignore")
}
