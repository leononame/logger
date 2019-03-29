package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
)

func ltolr(level Level) logrus.Level {
	switch level {
	case DebugLevel:
		return logrus.DebugLevel
	case InfoLevel:
		return logrus.InfoLevel
	case WarnLevel:
		return logrus.WarnLevel
	case ErrorLevel:
		return logrus.ErrorLevel
	case FatalLevel:
		return logrus.FatalLevel
	case PanicLevel:
		return logrus.PanicLevel
	}
	panic(fmt.Sprintf("Can't map level %d to logrus level", level))
}

func newLogrus(w io.Writer, lvl Level) Logger {
	l := logrus.New()
	l.SetOutput(w)
	l.SetLevel(ltolr(lvl))
	return &lLog{l}
}

type lLog struct {
	writer logrus.FieldLogger
}

func (l *lLog) WithField(key, value string) Logger {
	writer := l.writer.WithField(key, value)
	return &lLog{writer: writer}
}

func (l *lLog) Level(lvl Level) Entry {
	switch lvl {
	case DebugLevel:
		return l.Debug()
	case InfoLevel:
		return l.Info()
	case WarnLevel:
		return l.Warn()
	case ErrorLevel:
		return l.Error()
	case FatalLevel:
		return l.Fatal()
	case PanicLevel:
		return l.Panic()
	default:
		return l.Info()
	}
}

func (l *lLog) Debug() Entry {
	return &lEntry{logrus.DebugLevel, l.writer.WithField("time", time.Now())}
}
func (l *lLog) Info() Entry {
	return &lEntry{logrus.InfoLevel, l.writer.WithField("time", time.Now())}
}
func (l *lLog) Warn() Entry {
	return &lEntry{logrus.WarnLevel, l.writer.WithField("time", time.Now())}
}
func (l *lLog) Error() Entry {
	return &lEntry{logrus.ErrorLevel, l.writer.WithField("time", time.Now())}
}
func (l *lLog) Fatal() Entry {
	return &lEntry{logrus.FatalLevel, l.writer.WithField("time", time.Now())}
}
func (l *lLog) Panic() Entry {
	return &lEntry{logrus.PanicLevel, l.writer.WithField("time", time.Now())}
}

type lEntry struct {
	level logrus.Level
	entry *logrus.Entry
}

func (l *lEntry) Flush(msg string) {
	l.entry.Logln(l.level, msg)
	if l.level == logrus.FatalLevel {
		os.Exit(1)
	}
}

func (l *lEntry) AddFields(fs map[string]interface{}) Entry {
	l.entry = l.entry.WithFields(fs)
	return l
}

func (l *lEntry) AddErr(err error) Entry {
	msg := err.Error()
	st := errors.ErrorStack(err)
	l.entry = l.entry.WithField("err", msg)
	l.entry = l.entry.WithField("err_stack", st)
	return l
}

func (l *lEntry) AddError(key string, val error) Entry {
	msg := val.Error()
	st := errors.ErrorStack(val)
	l.entry = l.entry.WithField(key, msg)
	l.entry = l.entry.WithField(key+"_stack", st)
	return l
}

func (l *lEntry) AddBool(key string, val bool) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

func (l *lEntry) AddInt(key string, val int) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

func (l *lEntry) AddStr(key string, val string) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

func (l *lEntry) AddTime(key string, val time.Time) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

func (l *lEntry) AddDur(key string, val time.Duration) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}

func (l *lEntry) AddAny(key string, val interface{}) Entry {
	l.entry = l.entry.WithField(key, val)
	return l
}
