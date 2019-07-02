package logger

import "log"

// Logger interface should have Printf func
type Logger interface {
	Printf(format string, items ...interface{})
}

// NoopLogger does not log anything.
type NoopLogger struct{}

// Printf does nothing.
func (l NoopLogger) Printf(format string, items ...interface{}) {
	log.Printf(format, items...)
}
