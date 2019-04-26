package logger_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/juju/errors"

	"github.com/leononame/logger"
)

type l struct {
	buff *bytes.Buffer
	w    io.Writer
	lvl  logger.Level
}

func (l *l) WithField(key, value string) logger.Logger {
	_, _ = fmt.Fprint(l.buff, key)
	_, _ = fmt.Fprint(l.buff, value)
	return l
}

func (l *l) Level(lvl logger.Level) logger.Entry {
	if lvl > logger.DebugLevel || lvl < logger.PanicLevel {
		lvl = logger.InfoLevel
	}
	return &e{w: l.w, lvl: lvl, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Debug() logger.Entry {
	return &e{w: l.w, lvl: logger.DebugLevel, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Info() logger.Entry {
	return &e{w: l.w, lvl: logger.InfoLevel, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Warn() logger.Entry {
	return &e{w: l.w, lvl: logger.WarnLevel, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Error() logger.Entry {
	return &e{w: l.w, lvl: logger.ErrorLevel, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Fatal() logger.Entry {
	return &e{w: l.w, lvl: logger.FatalLevel, loglvl: l.lvl, buff: l.buff}
}

func (l *l) Panic() logger.Entry {
	return &e{w: l.w, lvl: logger.PanicLevel, loglvl: l.lvl, buff: l.buff}
}

type e struct {
	buff   *bytes.Buffer
	w      io.Writer
	lvl    logger.Level
	loglvl logger.Level
}

func (e *e) Flush(msg string) {
	if e.lvl <= e.loglvl {
		_, _ = e.w.Write(e.buff.Bytes())
		_, _ = e.w.Write([]byte(msg))
	}
	if e.lvl == logger.PanicLevel {
		panic(msg)
	}
	if e.lvl == logger.FatalLevel {
		os.Exit(1)
	}
}

func (e *e) AddFields(d map[string]interface{}) logger.Entry {
	for k, v := range d {
		_, _ = fmt.Fprint(e.buff, k)
		_, _ = fmt.Fprint(e.buff, v)
	}
	return e
}

func (e *e) AddErr(err error) logger.Entry {
	_, _ = fmt.Fprint(e.buff, "err")
	_, _ = fmt.Fprint(e.buff, err)
	_, _ = fmt.Fprint(e.buff, "err_stack")
	_, _ = fmt.Fprint(e.buff, errors.ErrorStack(err))
	return e
}

func (e *e) AddError(key string, val error) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	_, _ = fmt.Fprint(e.buff, key+"_stack")
	_, _ = fmt.Fprint(e.buff, errors.ErrorStack(val))
	return e
}

func (e *e) AddBool(key string, val bool) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}

func (e *e) AddInt(key string, val int) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}

func (e *e) AddStr(key string, val string) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}

func (e *e) AddTime(key string, val time.Time) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}

func (e *e) AddDur(key string, val time.Duration) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}

func (e *e) AddAny(key string, val interface{}) logger.Entry {
	_, _ = fmt.Fprint(e.buff, key)
	_, _ = fmt.Fprint(e.buff, val)
	return e
}
