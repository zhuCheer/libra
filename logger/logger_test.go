package logger

import "testing"

func TestLogger(t *testing.T) {
	var testB interface{} = &NoopLogger{}
	_, ok := testB.(Logger)

	if ok == false {
		t.Error("LoggerInterface implemention have an error #1")
	}
}

func TestPrintf(t *testing.T) {
	var logger = &NoopLogger{}

	logger.Printf("test logger %s", "string")
}
