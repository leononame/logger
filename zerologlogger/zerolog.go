package zerologlogger

import (
	"fmt"
	"io"
	"time"

	"github.com/leononame/logger"

	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

func ltoz(level logger.Level) zerolog.Level {
	switch level {
	case logger.DebugLevel:
		return zerolog.DebugLevel
	case logger.InfoLevel:
		return zerolog.InfoLevel
	case logger.WarnLevel:
		return zerolog.WarnLevel
	case logger.ErrorLevel:
		return zerolog.ErrorLevel
	case logger.FatalLevel:
		return zerolog.FatalLevel
	case logger.PanicLevel:
		return zerolog.PanicLevel
	}
	panic(fmt.Sprintf("Can't map level %d to zerolog level", level))
}

func New(w io.Writer, lvl logger.Level) logger.Logger {
	zerolog.TimeFieldFormat = ""
	l := zerolog.New(w).Level(ltoz(lvl)).With().Timestamp().Logger()
	return &zLog{&l}
}

type zLog struct {
	writer *zerolog.Logger
}

// WithField returns a new Logger that always logs the specified field
func (z *zLog) WithField(key, value string) logger.Logger {
	writer := z.writer.With().Str(key, value).Logger()
	return &zLog{writer: &writer}
}

// Level creates a new Entry with the specified Level
func (z *zLog) Level(lvl logger.Level) logger.Entry {
	switch lvl {
	case logger.DebugLevel:
		return z.Debug()
	case logger.InfoLevel:
		return z.Info()
	case logger.WarnLevel:
		return z.Warn()
	case logger.ErrorLevel:
		return z.Error()
	case logger.FatalLevel:
		return z.Fatal()
	case logger.PanicLevel:
		return z.Panic()
	default:
		return z.Info()
	}
}

// Debug creates a new Entry with level Debug
func (z *zLog) Debug() logger.Entry {
	return &zEntry{z.writer.Debug()}
}

// Info creates a new Entry with level Info
func (z *zLog) Info() logger.Entry {
	return &zEntry{z.writer.Info()}
}

// Warn creates a new Entry with level Warn
func (z *zLog) Warn() logger.Entry {
	return &zEntry{z.writer.Warn()}
}

// Error creates a new Entry with level Error
func (z *zLog) Error() logger.Entry {
	return &zEntry{z.writer.Error()}
}

// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
func (z *zLog) Fatal() logger.Entry {
	return &zEntry{z.writer.Fatal()}
}

// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
func (z *zLog) Panic() logger.Entry {
	return &zEntry{z.writer.Panic()}
}

type zEntry struct {
	entry *zerolog.Event
}

// Flush writes the entry as a single log statement. Optionally, a message can be added which will
// be included in the final log entry
func (z *zEntry) Flush(msg string) {
	z.entry.Msg(msg)
}

// AddFields adds a range of fields to the log statement
func (z *zEntry) AddFields(fs map[string]interface{}) logger.Entry {
	z.entry = z.entry.Fields(fs)
	return z
}

// AddErr adds an error to the log statement. The error will have the key "err". An error stack will be included
// under the key "err_stack"
func (z *zEntry) AddErr(err error) logger.Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	z.entry = z.entry.Str("err", msg)
	z.entry = z.entry.Str("err_stack", st)
	return z
}

// AddError adds an error to the log statement. An error stack will be included under the key "${key}_stack"
func (z *zEntry) AddError(key string, val error) logger.Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	z.entry = z.entry.Str(key, msg)
	z.entry = z.entry.Str(key+"_stack", st)
	z.entry = z.entry.AnErr(key, val)
	return z
}

// AddBool adds a bool value to the log statement.
func (z *zEntry) AddBool(key string, val bool) logger.Entry {
	z.entry = z.entry.Bool(key, val)
	return z
}

// AddInt adds an integer value to the log statement.
func (z *zEntry) AddInt(key string, val int) logger.Entry {
	z.entry = z.entry.Int(key, val)
	return z
}

// AddStr adds a string value to the log statement.
func (z *zEntry) AddStr(key string, val string) logger.Entry {
	z.entry = z.entry.Str(key, val)
	return z
}

// AddTime adds a time value to the log statement.
func (z *zEntry) AddTime(key string, val time.Time) logger.Entry {
	z.entry = z.entry.Time(key, val)
	return z
}

// AddDur adds a duration value to the log statement.
func (z *zEntry) AddDur(key string, val time.Duration) logger.Entry {
	z.entry = z.entry.Dur(key, val)
	return z
}

// AddAny adds any value to the log statement.
func (z *zEntry) AddAny(key string, val interface{}) logger.Entry {
	z.entry = z.entry.Interface(key, val)
	return z
}
