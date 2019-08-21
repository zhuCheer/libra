package logger

import (
	"log"
	"strings"
)

// Logger interface should have Printf func
type Logger interface {
	Printf(format string, items ...interface{})
}

// NoopLogger does not log anything.
type NoopLogger struct{}

const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	CRITICAL
)

var loggerLevel = DEBUG

// SetLevel set logger level
func (l NoopLogger) SetLevel(level string) {
	level = strings.ToUpper(level)

	switch level {
	case "DEBUG":
		loggerLevel = DEBUG
	case "INFO":
		loggerLevel = INFO
	case "WARN":
		loggerLevel = WARN
	case "ERROR":
		loggerLevel = ERROR
	case "CRITICAL":
		loggerLevel = CRITICAL
	default:
		loggerLevel = DEBUG
	}
}

// Printf print log.
func (l NoopLogger) Printf(format string, items ...interface{}) {
	log.Printf(format, items...)
}

// Debug print Debug level log.
func (l NoopLogger) Debug(format string, items ...interface{}) {
	if loggerLevel <= DEBUG {
		log.Printf("[DEBUG]"+format, items...)
	}
}

// Info print Info level.
func (l NoopLogger) Info(format string, items ...interface{}) {
	if loggerLevel <= INFO {
		log.Printf("[INFO]"+format, items...)
	}
}

// Warn print Warn level.
func (l NoopLogger) Warn(format string, items ...interface{}) {
	if loggerLevel <= WARN {
		log.Printf("[Warn]"+format, items...)
	}
}

// Error print Error level.
func (l NoopLogger) Error(format string, items ...interface{}) {
	if loggerLevel <= ERROR {
		log.Printf("[ERROR]"+format, items...)
	}
}

// Critical print Critical level.
func (l NoopLogger) Critical(format string, items ...interface{}) {
	if loggerLevel <= CRITICAL {
		log.Printf("[CRITICAL]"+format, items...)
	}
}
