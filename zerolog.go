package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/juju/errors"
	"github.com/rs/zerolog"
)

func ltoz(level Level) zerolog.Level {
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

func newZeroLog(w io.Writer, lvl Level) Logger {
	zerolog.TimeFieldFormat = ""
	l := zerolog.New(w).Level(ltoz(lvl)).With().Timestamp().Logger()
	return &zLog{&l}
}

type zLog struct {
	writer *zerolog.Logger
}

func (z *zLog) WithField(key, value string) Logger {
	writer := z.writer.With().Str(key, value).Logger()
	return &zLog{writer: &writer}
}

func (z *zLog) Level(lvl Level) Entry {
	switch lvl {
	case DebugLevel:
		return z.Debug()
	case InfoLevel:
		return z.Info()
	case WarnLevel:
		return z.Warn()
	case ErrorLevel:
		return z.Error()
	case FatalLevel:
		return z.Fatal()
	case PanicLevel:
		return z.Panic()
	default:
		return z.Info()
	}
}

func (z *zLog) Debug() Entry {
	return &zEntry{z.writer.Debug()}
}
func (z *zLog) Info() Entry {
	return &zEntry{z.writer.Info()}
}
func (z *zLog) Warn() Entry {
	return &zEntry{z.writer.Warn()}
}
func (z *zLog) Error() Entry {
	return &zEntry{z.writer.Error()}
}
func (z *zLog) Fatal() Entry {
	return &zEntry{z.writer.Fatal()}
}
func (z *zLog) Panic() Entry {
	return &zEntry{z.writer.Panic()}
}

type zEntry struct {
	entry *zerolog.Event
}

func (z *zEntry) Flush(msg string) {
	z.entry.Msg(msg)
}

func (z *zEntry) AddFields(fs map[string]interface{}) Entry {
	z.entry = z.entry.Fields(fs)
	return z
}

func (z *zEntry) AddErr(err error) Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	z.entry = z.entry.Str("err", msg)
	z.entry = z.entry.Str("err_stack", st)
	return z
}

func (z *zEntry) AddError(key string, val error) Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	z.entry = z.entry.Str(key, msg)
	z.entry = z.entry.Str(key+"_stack", st)
	z.entry = z.entry.AnErr(key, val)
	return z
}

func (z *zEntry) AddBool(key string, val bool) Entry {
	z.entry = z.entry.Bool(key, val)
	return z
}

func (z *zEntry) AddInt(key string, val int) Entry {
	z.entry = z.entry.Int(key, val)
	return z
}

func (z *zEntry) AddStr(key string, val string) Entry {
	z.entry = z.entry.Str(key, val)
	return z
}

func (z *zEntry) AddTime(key string, val time.Time) Entry {
	z.entry = z.entry.Time(key, val)
	return z
}

func (z *zEntry) AddDur(key string, val time.Duration) Entry {
	z.entry = z.entry.Dur(key, val)
	return z
}

func (z *zEntry) AddAny(key string, val interface{}) Entry {
	z.entry = z.entry.Interface(key, val)
	return z
}
