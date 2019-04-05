package logger_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/leononame/logger"
	"github.com/leononame/logger/logruslogger"
	"github.com/leononame/logger/zerologlogger"

	"github.com/leononame/logger/gelflogger"

	"github.com/juju/errors"
	. "github.com/leononame/logger"
	"github.com/stretchr/testify/assert"
)

const IncorrectLevel logger.Level = 1000

func tests() []struct {
	lvl logger.Level
	key string
	val interface{}
} {
	return []struct {
		lvl logger.Level
		key string
		val interface{}
	}{
		{logger.DebugLevel, "lvldebug", "true"},
		{logger.InfoLevel, "lvlinfo", "true"},
		{logger.WarnLevel, "lvlwarn", "true"},
		{logger.ErrorLevel, "lvlerr", "true"},
		{logger.PanicLevel, "lvlpanic", "true"},
	}
}

func multiLogger(lvl Level) (Logger, []strings.Builder) {
	sbs := make([]strings.Builder, 3)
	l0 := gelflogger.New(&sbs[0], lvl)
	l1 := zerologlogger.New(&sbs[1], lvl)
	l2 := logruslogger.New(&sbs[2], lvl)
	return NewMulti(l0, l1, l2), sbs
}

func TestMLog_WithField(t *testing.T) {
	l, sbs := multiLogger(DebugLevel)
	l.WithField("somekey", "someval")
	l.Debug().AddStr("otherkey", "otherval").Flush("message")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, "somekey", "Log message should contain key")
		assert.Contains(t, s, "someval", "Log message should contain value")
	}
}

func TestMLog_Level(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(DebugLevel)

		f := func() { l.Level(test.lvl).AddAny(test.key, test.val).Flush("Message") }
		if test.lvl == PanicLevel {
			assert.Panics(t, f, "Function should panic")
		} else {
			f()
		}
		for _, sb := range sbs {
			s := sb.String()
			assert.Contains(t, s, test.key, "Logger should print message")
		}
	}
}

func TestMLog_Level2(t *testing.T) {
	l, sbs := multiLogger(DebugLevel)
	l.Level(IncorrectLevel).AddAny("somekey", "someval").Flush("Message")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, "somekey", "Logger should print message")
		assert.Contains(t, s, "somekey", "Logger should print message")
	}
}

func TestMLog_Debug(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(test.lvl)
		l.Debug().AddAny(test.key, test.val).Flush("Additional Message")
		for _, sb := range sbs {
			s := sb.String()
			if test.lvl >= DebugLevel {
				msg := fmt.Sprintf("Logger with level %d should print Debug messages", test.lvl)
				assert.Contains(t, s, test.key, msg)
			} else {
				msg := fmt.Sprintf("Logger with level %d should not print Debug messages", test.lvl)
				assert.NotContains(t, s, test.key, msg)
			}
		}
	}
}

func TestMLog_Info(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(test.lvl)
		l.Info().AddAny(test.key, test.val).Flush("Additional Message")
		for _, sb := range sbs {
			s := sb.String()
			if test.lvl >= InfoLevel {
				msg := fmt.Sprintf("Logger with level %d should print Info messages", test.lvl)
				assert.Contains(t, s, test.key, msg)
			} else {
				msg := fmt.Sprintf("Logger with level %d should not print Info messages", test.lvl)
				assert.NotContains(t, s, test.key, msg)
			}
		}
	}
}

func TestMLog_Warn(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(test.lvl)
		l.Warn().AddAny(test.key, test.val).Flush("Additional Message")
		for _, sb := range sbs {
			s := sb.String()
			if test.lvl >= WarnLevel {
				msg := fmt.Sprintf("Logger with level %d should print Warn messages", test.lvl)
				assert.Contains(t, s, test.key, msg)
			} else {
				msg := fmt.Sprintf("Logger with level %d should not print Warn messages", test.lvl)
				assert.NotContains(t, s, test.key, msg)
			}
		}
	}
}

func TestMLog_Error(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(test.lvl)
		l.Error().AddAny(test.key, test.val).Flush("Additional Message")
		for _, sb := range sbs {
			s := sb.String()
			if test.lvl >= ErrorLevel {
				msg := fmt.Sprintf("Logger with level %d should print Error messages", test.lvl)
				assert.Contains(t, s, test.key, msg)
			} else {
				msg := fmt.Sprintf("Logger with level %d should not print Error messages", test.lvl)
				assert.NotContains(t, s, test.key, msg)
			}
		}
	}
}

func TestMLog_Panic(t *testing.T) {
	for _, test := range tests() {
		l, sbs := multiLogger(test.lvl)
		e := l.Panic().AddAny(test.key, test.val)
		f := func() {
			e.Flush("Additional Message")
		}
		assert.Panics(t, f, "Call to Panic level should panic")
		for _, sb := range sbs {
			s := sb.String()
			if test.lvl >= PanicLevel {
				msg := fmt.Sprintf("Logger with level %d should print Panic messages", test.lvl)
				assert.Contains(t, s, test.key, msg)
			} else {
				msg := fmt.Sprintf("Logger with level %d should not print Panic messages", test.lvl)
				assert.NotContains(t, s, test.key, msg)
			}
		}
	}
}

func TestMEntry_AddBool(t *testing.T) {
	key := "boolkey"
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddBool(key, true).Flush("")
	for _, sb := range sbs {
		assert.Contains(t, sb.String(), key, "Message should contain key")
	}
}

func TestMEntry_AddDur(t *testing.T) {
	key := "durkey"
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddDur(key, time.Since(time.Now())).Flush("")
	for _, sb := range sbs {
		assert.Contains(t, sb.String(), key, "Message should contain key")
	}
}

func TestMEntry_AddAny(t *testing.T) {
	key := "anykey"
	val := "valval"
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddAny(key, val).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, key, "Message should contain key")
		assert.Contains(t, s, val, "Message should contain value")
	}
}

func TestMEntry_AddErr(t *testing.T) {
	l, sbs := multiLogger(DebugLevel)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddErr(err2).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, "err_stack", "Message should contain error stack")
		assert.Contains(t, s, "asd", "Message should contain error mesage")
		assert.Contains(t, s, "other err", "Message should contain other error mesage")
	}
}

func TestMEntry_AddError(t *testing.T) {
	key := "errkey"
	l, sbs := multiLogger(DebugLevel)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddError(key, err2).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, key, "Message should contain key")
		assert.Contains(t, s, key+"_stack", "Message should contain key stack")
		assert.Contains(t, s, "asd", "Message should contain error mesage")
		assert.Contains(t, s, "other err", "Message should contain other error mesage")
	}
}

func TestMEntry_AddFields(t *testing.T) {
	data := map[string]interface{}{
		// avoid time as value because we don't control formatting necessarily
		"key1":      "strval",
		"other_key": false,
		"third key": 9999,
	}
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddFields(data).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		for k, v := range data {
			assert.Contains(t, s, k, "Log message should contain key "+k)
			assert.Contains(t, s, fmt.Sprint(v), "Log message should contain value")
		}
	}
}

func TestMEntry_AddInt(t *testing.T) {
	key := "intkey"
	val := 1990123
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddInt(key, val).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, key, "Message should contain key")
		assert.Contains(t, s, fmt.Sprint(val), "Message should contain value")
	}
}

func TestMEntry_AddStr(t *testing.T) {
	key := "strkey"
	val := "thisisavalue"
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddStr(key, val).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, key, "Message should contain key")
		assert.Contains(t, s, val, "Message should contain value")

	}
}

func TestMEntry_AddTime(t *testing.T) {
	key := "timekey"
	val := time.Now()
	l, sbs := multiLogger(DebugLevel)
	l.Info().AddTime(key, val).Flush("")
	for _, sb := range sbs {
		s := sb.String()
		assert.Contains(t, s, key, "Message should contain key")
	}
}
