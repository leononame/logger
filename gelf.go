package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

func ltog(level Level) zerolog.Level {
	switch level {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	case PanicLevel:
		return zerolog.PanicLevel
	}
	panic(fmt.Sprintf("Can't map level %d to zerolog level", level))
}

func newGelfLog(w io.Writer, lvl Level) Logger {
	l := zerolog.New(w).Level(ltog(lvl))
	return &gLog{&l, lvl}
}

type gLog struct {
	writer *zerolog.Logger
	level  Level
}

func (g *gLog) WithField(key, value string) Logger {
	writer := g.writer.With().Str("_"+key, value).Logger()
	return &gLog{writer: &writer, level: g.level}
}

func (g *gLog) Level(lvl Level) Entry {
	switch lvl {
	case DebugLevel:
		return g.Debug()
	case InfoLevel:
		return g.Info()
	case WarnLevel:
		return g.Warn()
	case ErrorLevel:
		return g.Error()
	case FatalLevel:
		return g.Fatal()
	case PanicLevel:
		return g.Panic()
	default:
		return g.Info()
	}
}

func (g *gLog) Debug() Entry {
	l := g.writer.With().Int("level", int(DebugLevel)).Logger()
	var e = l.Log()
	if g.level < DebugLevel {
		e = nil
	}
	return &gEntry{e, DebugLevel}
}
func (g *gLog) Info() Entry {
	l := g.writer.With().Int("level", int(InfoLevel)).Logger()
	var e = l.Log()
	if g.level < InfoLevel {
		e = nil
	}
	return &gEntry{e, InfoLevel}
}
func (g *gLog) Warn() Entry {
	l := g.writer.With().Int("level", int(WarnLevel)).Logger()
	var e = l.Log()
	if g.level < WarnLevel {
		e = nil
	}
	return &gEntry{e, WarnLevel}
}
func (g *gLog) Error() Entry {
	l := g.writer.With().Int("level", int(ErrorLevel)).Logger()
	var e = l.Log()
	if g.level < ErrorLevel {
		e = nil
	}
	return &gEntry{e, ErrorLevel}
}
func (g *gLog) Fatal() Entry {
	l := g.writer.With().Int("level", int(FatalLevel)).Logger()
	var e = l.Log()
	if g.level < FatalLevel {
		e = nil
	}
	return &gEntry{e, FatalLevel}
}
func (g *gLog) Panic() Entry {
	l := g.writer.With().Int("level", int(PanicLevel)).Logger()
	var e = l.Log()
	if g.level < PanicLevel {
		e = nil
	}
	return &gEntry{e, PanicLevel}
}

type gEntry struct {
	entry *zerolog.Event
	lvl   Level
}

func (g *gEntry) Flush(msg string) {
	g.entry.Int64("timestamp", time.Now().Unix())
	g.entry.Str("version", "1.1")
	g.entry.Str("short_message", msg)
	// This skips a message in zerolog
	g.entry.Msg("")
	if g.lvl == PanicLevel {
		panic("logger called at panic level with message: " + msg)
	} else if g.lvl == FatalLevel {
		os.Exit(1)
	}
}

func (g *gEntry) AddFields(fs map[string]interface{}) Entry {
	for k, v := range fs {
		g.entry = g.entry.Interface("_"+k, v)
	}
	return g
}

func (g *gEntry) AddErr(err error) Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	g.entry = g.entry.Str("_err", msg)
	g.entry = g.entry.Str("_err_stack", st)
	return g
}

func (g *gEntry) AddError(key string, val error) Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	g.entry = g.entry.Str("_"+key, msg)
	g.entry = g.entry.Str("_"+key+"_stack", st)
	g.entry = g.entry.AnErr("_"+key, val)
	return g
}

func (g *gEntry) AddBool(key string, val bool) Entry {
	g.entry = g.entry.Bool("_"+key, val)
	return g
}

func (g *gEntry) AddInt(key string, val int) Entry {
	g.entry = g.entry.Int("_"+key, val)
	return g
}

func (g *gEntry) AddStr(key string, val string) Entry {
	g.entry = g.entry.Str("_"+key, val)
	return g
}

func (g *gEntry) AddTime(key string, val time.Time) Entry {
	g.entry = g.entry.Time("_"+key, val)
	return g
}

func (g *gEntry) AddDur(key string, val time.Duration) Entry {
	g.entry = g.entry.Dur("_"+key, val)
	return g
}

func (g *gEntry) AddAny(key string, val interface{}) Entry {
	g.entry = g.entry.Interface("_"+key, val)
	return g
}
