package gelflogger

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/leononame/logger"

	"github.com/juju/errors"
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

func TestGLog_WithField(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel).WithField("somekey", "someval")
	l.Debug().AddStr("otherkey", "otherval").Flush("message")
	s := sb.String()
	assert.Contains(t, s, "_somekey", "Log message should contain key")
	assert.Contains(t, s, "someval", "Log message should contain value")
}

func TestGLog_Level(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)

		f := func() { l.Level(test.lvl).AddAny(test.key, test.val).Flush("Message") }
		if test.lvl == logger.PanicLevel {
			assert.Panics(t, f, "Function should panic")
		} else {
			f()
		}
		s := sb.String()
		assert.Contains(t, s, "_"+test.key, "Logger should print message")
	}
}

func TestGLog_Level2(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, logger.InfoLevel)
	l.Level(IncorrectLevel).AddAny("somekey", "someval").Flush("Message")
	s := sb.String()
	assert.Contains(t, s, "_somekey", "Logger should print message")
	assert.Contains(t, s, `"level":6`, "Logger should print at info level")
}

func TestGLog_Debug(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)
		l.Debug().AddAny(test.key, test.val).Flush("Additional Message")
		s := sb.String()
		if test.lvl >= logger.DebugLevel {
			msg := fmt.Sprintf("Logger with level %d should print Debug messages", test.lvl)
			assert.Contains(t, s, "_"+test.key, msg)
		} else {
			msg := fmt.Sprintf("Logger with level %d should not print Debug messages", test.lvl)
			assert.NotContains(t, s, test.key, msg)
		}
	}
}

func TestGLog_Info(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)
		l.Info().AddAny(test.key, test.val).Flush("Additional Message")
		s := sb.String()
		if test.lvl >= logger.InfoLevel {
			msg := fmt.Sprintf("Logger with level %d should print Info messages", test.lvl)
			assert.Contains(t, s, "_"+test.key, msg)
		} else {
			msg := fmt.Sprintf("Logger with level %d should not print Info messages", test.lvl)
			assert.NotContains(t, s, test.key, msg)
		}
	}
}

func TestGLog_Warn(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)
		l.Warn().AddAny(test.key, test.val).Flush("Additional Message")
		s := sb.String()
		if test.lvl >= logger.WarnLevel {
			msg := fmt.Sprintf("Logger with level %d should print Warn messages", test.lvl)
			assert.Contains(t, s, "_"+test.key, msg)
		} else {
			msg := fmt.Sprintf("Logger with level %d should not print Warn messages", test.lvl)
			assert.NotContains(t, s, test.key, msg)
		}
	}
}

func TestGLog_Error(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)
		l.Error().AddAny(test.key, test.val).Flush("Additional Message")
		s := sb.String()
		if test.lvl >= logger.ErrorLevel {
			msg := fmt.Sprintf("Logger with level %d should print Error messages", test.lvl)
			assert.Contains(t, s, "_"+test.key, msg)
		} else {
			msg := fmt.Sprintf("Logger with level %d should not print Error messages", test.lvl)
			assert.NotContains(t, s, test.key, msg)
		}
	}
}

func TestGLog_Panic(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl)
		e := l.Panic().AddAny(test.key, test.val)
		f := func() {
			e.Flush("Additional Message")
		}
		assert.Panics(t, f, "Call to Panic level should panic")
		s := sb.String()
		if test.lvl >= logger.PanicLevel {
			msg := fmt.Sprintf("Logger with level %d should print Panic messages", test.lvl)
			assert.Contains(t, s, "_"+test.key, msg)
		} else {
			msg := fmt.Sprintf("Logger with level %d should not print Panic messages", test.lvl)
			assert.NotContains(t, s, test.key, msg)
		}
	}
}

func TestGEntry_AddBool(t *testing.T) {
	key := "boolkey"
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddBool(key, true).Flush("")
	assert.Contains(t, sb.String(), "_"+key, "Message should contain key")
}

func TestGEntry_AddDur(t *testing.T) {
	key := "durkey"
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddDur(key, time.Since(time.Now())).Flush("")
	assert.Contains(t, sb.String(), "_"+key, "Message should contain key")
}

func TestGEntry_AddAny(t *testing.T) {
	key := "anykey"
	val := "valval"
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddAny(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_"+key, "Message should contain key")
	assert.Contains(t, s, val, "Message should contain value")
}

func TestGEntry_AddErr(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddErr(err2).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_err_stack", "Message should contain error stack")
	assert.Contains(t, s, "asd", "Message should contain error mesage")
	assert.Contains(t, s, "other err", "Message should contain other error mesage")
}

func TestGEntry_AddError(t *testing.T) {
	key := "errkey"
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddError(key, err2).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_"+key, "Message should contain key")
	assert.Contains(t, s, "_"+key+"_stack", "Message should contain key stack")
	assert.Contains(t, s, "asd", "Message should contain error mesage")
	assert.Contains(t, s, "other err", "Message should contain other error mesage")
}

func TestGEntry_AddFields(t *testing.T) {
	data := map[string]interface{}{
		// avoid time as value because we don't control formatting necessarily
		"key1":      "strval",
		"other_key": false,
		"third key": 9999,
	}
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddFields(data).Flush("")
	s := sb.String()
	for k, v := range data {
		assert.Contains(t, s, "_"+k, "Log message should contain key "+k)
		assert.Contains(t, s, fmt.Sprint(v), "Log message should contain value")
	}
}

func TestGEntry_AddInt(t *testing.T) {
	key := "intkey"
	val := 1990123
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddInt(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_"+key, "Message should contain key")
	assert.Contains(t, s, fmt.Sprint(val), "Message should contain value")
}

func TestGEntry_AddStr(t *testing.T) {
	key := "strkey"
	val := "thisisavalue"
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddStr(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_"+key, "Message should contain key")
	assert.Contains(t, s, val, "Message should contain value")

}

func TestGEntry_AddTime(t *testing.T) {
	key := "timekey"
	val := time.Now()
	var sb strings.Builder
	l := New(&sb, logger.DebugLevel)
	l.Info().AddTime(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, "_"+key, "Message should contain key")
}
