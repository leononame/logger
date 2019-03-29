package logger

import (
	"io"
	"time"
)

type Implementation int

const (
	// ZeroLogBackend defines zerolog as the actual log implementation
	ZeroLogBackend Implementation = iota
	// LogrusBackend defines logrus as the actual log implementation
	LogrusBackend
)

type Level int

// Levels have the same value as syslog, hence 5 is skipped
const (
	// DebugLevel defines debug log level.
	DebugLevel Level = 7
	// InfoLevel defines info log level.
	InfoLevel = 6
	// WarnLevel defines warn log level.
	WarnLevel = 4
	// ErrorLevel defines error log level.
	ErrorLevel = 3
	// FatalLevel defines fatal log level.
	FatalLevel = 2
	// PanicLevel defines panic log level.
	PanicLevel = 1
)

// New returns a logger. The logger will write to the writer specified and will use the log backend specified
func New(w io.Writer, lvl Level, impl Implementation) Logger {
	var l Logger
	// only one implementation, always go to default case
	switch impl {
	case LogrusBackend:
		l = newLogrus(w, lvl)
	case ZeroLogBackend:
		fallthrough
	default:
		l = newZeroLog(w, lvl)
	}
	return l
}

// Logger is an standard interface for logging so that different log implementations can be wrapped around.
// The API is heavily influenced by the zerolog API for structured JSON logging
type Logger interface {
	// WithField returns a new Logger that always logs the specified field
	WithField(key, value string) Logger
	// Level creates a new Entry with the specified Level
	Level(Level) Entry
	// Debug creates a new Entry with level Debug
	Debug() Entry
	// Info creates a new Entry with level Info
	Info() Entry
	// Warn creates a new Entry with level Warn
	Warn() Entry
	// Error creates a new Entry with level Error
	Error() Entry
	// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
	Fatal() Entry
	// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
	Panic() Entry
}

// Entry is an interface for a log entry. A single entry always has defined a log level. Custom fields can be
// added. An entry will never be written to the log unless Flush is called.
type Entry interface {
	// Flush writes the entry as a single log statement. Optionally, a message can be added which will
	// be included in the final log entry
	Flush(string)

	// Add a range of fields to the log statement
	AddFields(map[string]interface{}) Entry
	// Add an error to the log statement. The error will have the key "err". An error stack will be included
	// under the key "err_stack"
	AddErr(err error) Entry
	// Add an error to the log statement. An error stack will be included under the key "${key}_stack"
	AddError(key string, val error) Entry
	// Add a bool value to the log statement.
	AddBool(key string, val bool) Entry
	// Add an integer value to the log statement.
	AddInt(key string, val int) Entry
	// Add a string value to the log statement.
	AddStr(key string, val string) Entry
	// Add a time value to the log statement.
	AddTime(key string, val time.Time) Entry
	// Add a duration value to the log statement.
	AddDur(key string, val time.Duration) Entry
	// Add any value to the log statement.
	AddAny(key string, val interface{}) Entry
}
