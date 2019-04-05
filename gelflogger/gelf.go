package gelflogger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/leononame/logger"

	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

func ltog(level logger.Level) zerolog.Level {
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
	l := zerolog.New(w).Level(ltog(lvl))
	return &gLog{&l, lvl}
}

type gLog struct {
	writer *zerolog.Logger
	level  logger.Level
}

// WithField returns a new Logger that always logs the specified field
func (g *gLog) WithField(key, value string) logger.Logger {
	writer := g.writer.With().Str("_"+key, value).Logger()
	return &gLog{writer: &writer, level: g.level}
}

// Level creates a new Entry with the specified Level
func (g *gLog) Level(lvl logger.Level) logger.Entry {
	switch lvl {
	case logger.DebugLevel:
		return g.Debug()
	case logger.InfoLevel:
		return g.Info()
	case logger.WarnLevel:
		return g.Warn()
	case logger.ErrorLevel:
		return g.Error()
	case logger.FatalLevel:
		return g.Fatal()
	case logger.PanicLevel:
		return g.Panic()
	default:
		return g.Info()
	}
}

// Debug creates a new Entry with level Debug
func (g *gLog) Debug() logger.Entry {
	l := g.writer.With().Int("level", int(logger.DebugLevel)).Logger()
	var e = l.Log()
	if g.level < logger.DebugLevel {
		e = nil
	}
	return &gEntry{e, logger.DebugLevel}
}

// Info creates a new Entry with level Info
func (g *gLog) Info() logger.Entry {
	l := g.writer.With().Int("level", int(logger.InfoLevel)).Logger()
	var e = l.Log()
	if g.level < logger.InfoLevel {
		e = nil
	}
	return &gEntry{e, logger.InfoLevel}
}

// Warn creates a new Entry with level Warn
func (g *gLog) Warn() logger.Entry {
	l := g.writer.With().Int("level", int(logger.WarnLevel)).Logger()
	var e = l.Log()
	if g.level < logger.WarnLevel {
		e = nil
	}
	return &gEntry{e, logger.WarnLevel}
}

// Error creates a new Entry with level Error
func (g *gLog) Error() logger.Entry {
	l := g.writer.With().Int("level", int(logger.ErrorLevel)).Logger()
	var e = l.Log()
	if g.level < logger.ErrorLevel {
		e = nil
	}
	return &gEntry{e, logger.ErrorLevel}
}

// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
func (g *gLog) Fatal() logger.Entry {
	l := g.writer.With().Int("level", int(logger.FatalLevel)).Logger()
	var e = l.Log()
	if g.level < logger.FatalLevel {
		e = nil
	}
	return &gEntry{e, logger.FatalLevel}
}

// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
func (g *gLog) Panic() logger.Entry {
	l := g.writer.With().Int("level", int(logger.PanicLevel)).Logger()
	var e = l.Log()
	if g.level < logger.PanicLevel {
		e = nil
	}
	return &gEntry{e, logger.PanicLevel}
}

type gEntry struct {
	entry *zerolog.Event
	lvl   logger.Level
}

// Flush writes the entry as a single log statement. Optionally, a message can be added which will
// be included in the final log entry
func (g *gEntry) Flush(msg string) {
	g.entry.Int64("timestamp", time.Now().Unix())
	g.entry.Str("version", "1.1")
	g.entry.Str("short_message", msg)
	// This skips a message in zerolog
	g.entry.Msg("")
	if g.lvl == logger.PanicLevel {
		panic("logger called at panic level with message: " + msg)
	} else if g.lvl == logger.FatalLevel {
		os.Exit(1)
	}
}

// AddFields adds a range of fields to the log statement
func (g *gEntry) AddFields(fs map[string]interface{}) logger.Entry {
	for k, v := range fs {
		g.entry = g.entry.Interface("_"+k, v)
	}
	return g
}

// AddErr adds an error to the log statement. The error will have the key "err". An error stack will be included
// under the key "err_stack"
func (g *gEntry) AddErr(err error) logger.Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	g.entry = g.entry.Str("_err", msg)
	g.entry = g.entry.Str("_err_stack", st)
	return g
}

// AddError adds an error to the log statement. An error stack will be included under the key "${key}_stack"
func (g *gEntry) AddError(key string, val error) logger.Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	g.entry = g.entry.Str("_"+key, msg)
	g.entry = g.entry.Str("_"+key+"_stack", st)
	g.entry = g.entry.AnErr("_"+key, val)
	return g
}

// AddBool adds a bool value to the log statement.
func (g *gEntry) AddBool(key string, val bool) logger.Entry {
	g.entry = g.entry.Bool("_"+key, val)
	return g
}

// AddInt adds an integer value to the log statement.
func (g *gEntry) AddInt(key string, val int) logger.Entry {
	g.entry = g.entry.Int("_"+key, val)
	return g
}

// AddStr adds a string value to the log statement.
func (g *gEntry) AddStr(key string, val string) logger.Entry {
	g.entry = g.entry.Str("_"+key, val)
	return g
}

// AddTime adds a time value to the log statement.
func (g *gEntry) AddTime(key string, val time.Time) logger.Entry {
	g.entry = g.entry.Time("_"+key, val)
	return g
}

// AddDur adds a duration value to the log statement.
func (g *gEntry) AddDur(key string, val time.Duration) logger.Entry {
	g.entry = g.entry.Dur("_"+key, val)
	return g
}

// AddAny adds any value to the log statement.
func (g *gEntry) AddAny(key string, val interface{}) logger.Entry {
	g.entry = g.entry.Interface("_"+key, val)
	return g
}
