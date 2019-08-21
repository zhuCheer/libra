package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	var testB interface{} = &NoopLogger{}
	_, ok := testB.(Logger)

	if ok == false {
		t.Error("LoggerInterface implemention have an error #1")
	}
}

func TestSetLevel(t *testing.T) {

	var logger = &NoopLogger{}
	logger.SetLevel("info")
	if loggerLevel != 1 {
		t.Error("logger SetLevel have an error #1")
	}

	logger.SetLevel("Warn")
	if loggerLevel != 2 {
		t.Error("logger SetLevel have an error #2")
	}

	logger.SetLevel("ERROR")
	if loggerLevel != 3 {
		t.Error("logger SetLevel have an error #3")
	}

	logger.SetLevel("debug")
	if loggerLevel != 0 {
		t.Error("logger SetLevel have an error #4")
	}

	logger.SetLevel("CRITICAL")
	if loggerLevel != 4 {
		t.Error("logger SetLevel have an error #5")
	}

	logger.SetLevel("xxx")
	if loggerLevel != 0 {
		t.Error("logger SetLevel have an error #6")
	}
}

func TestPrintf(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Printf("test logger %s", "string")
}

func TestDebug(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Debug("test Debug logger %s", "string")
}

func TestInfo(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Info("test info logger %s", "string")
}

func TestWarn(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Warn("test Warn logger %s", "string")
}

func TestError(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Error("test error logger %s", "string")
}

func TestCritical(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Critical("test Critical logger %s", "string")
}
