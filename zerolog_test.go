package logger

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/juju/errors"
	"github.com/stretchr/testify/assert"
)

func TestZLog_Level(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		f := func() {l.Level(test.lvl).AddAny(test.key, test.val).Flush("Message")}
		if test.lvl  == PanicLevel {
			assert.Panics(t, f, "Function should panic")
		} else {
			f()
		}
		s := sb.String()
		assert.Contains(t, s, test.key, "Logger should print message")
	}
}

func TestZLog_Level2(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, InfoLevel, ZeroLogBackend)
	l.Level(IncorrectLevel).AddAny("somekey", "someval").Flush("Message")
	s := sb.String()
	assert.Contains(t, s, "somekey", "Logger should print message")
	assert.Contains(t, s, "info", "Logger should print at info level")
}

func TestZLog_WithField(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend).WithField("somekey", "someval")
	l.Debug().AddStr("otherkey", "otherval").Flush("message")
	s := sb.String()
	assert.Contains(t, s, "somekey", "Log message should contain key")
	assert.Contains(t, s, "someval", "Log message should contain value")
}

func TestZLog_Debug(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		l.Debug().AddAny(test.key, test.val).Flush("Additional Message")
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

func TestZLog_Info(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		l.Info().AddAny(test.key, test.val).Flush("Additional Message")
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

func TestZLog_Warn(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		l.Warn().AddAny(test.key, test.val).Flush("Additional Message")
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

func TestZLog_Error(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		l.Error().AddAny(test.key, test.val).Flush("Additional Message")
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

func TestZLog_Panic(t *testing.T) {
	for _, test := range tests() {
		var sb strings.Builder
		l := New(&sb, test.lvl, ZeroLogBackend)
		e := l.Panic().AddAny(test.key, test.val)
		f := func() {
			e.Flush("Additional Message")
		}
		assert.Panics(t, f, "Call to Panic level should panic")
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

func TestZEntry_AddBool(t *testing.T) {
	key := "boolkey"
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddBool(key, true).Flush("")
	assert.Contains(t, sb.String(), key, "Message should contain key")
}

func TestZEntry_AddDur(t *testing.T) {
	key := "durkey"
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddDur(key, time.Since(time.Now())).Flush("")
	assert.Contains(t, sb.String(), key, "Message should contain key")
}

func TestZEntry_AddAny(t *testing.T) {
	key := "anykey"
	val := "valval"
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddAny(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, key, "Message should contain key")
	assert.Contains(t, s, val, "Message should contain value")
}

func TestZEntry_AddErr(t *testing.T) {
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddErr(err2).Flush("")
	s := sb.String()
	assert.Contains(t, s, "err_stack", "Message should contain error stack")
	assert.Contains(t, s, "asd", "Message should contain error mesage")
	assert.Contains(t, s, "other err", "Message should contain other error mesage")
}

func TestZEntry_AddError(t *testing.T) {
	key := "errkey"
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	err := errors.New("asd")
	err2 := errors.Annotate(err, "other err")
	l.Info().AddError(key, err2).Flush("")
	s := sb.String()
	assert.Contains(t, s, key, "Message should contain key")
	assert.Contains(t, s, key+"_stack", "Message should contain key stack")
	assert.Contains(t, s, "asd", "Message should contain error mesage")
	assert.Contains(t, s, "other err", "Message should contain other error mesage")
}

func TestZEntry_AddFields(t *testing.T) {
	data := map[string]interface{}{
		// avoid time as value because we don't control formatting necessarily
		"key1":      "strval",
		"other_key": false,
		"third key": 9999,
	}
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddFields(data).Flush("")
	s := sb.String()
	for k, v := range data {
		assert.Contains(t, s, k, "Log message should contain key "+k)
		assert.Contains(t, s, fmt.Sprint(v), "Log message should contain value")
	}
}

func TestZEntry_AddInt(t *testing.T) {
	key := "intkey"
	val := 1990123
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddInt(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, key, "Message should contain key")
	assert.Contains(t, s, fmt.Sprint(val), "Message should contain value")
}

func TestZEntry_AddStr(t *testing.T) {
	key := "strkey"
	val := "thisisavalue"
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddStr(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, key, "Message should contain key")
	assert.Contains(t, s, val, "Message should contain value")

}

func TestZEntry_AddTime(t *testing.T) {
	key := "timekey"
	val := time.Now()
	var sb strings.Builder
	l := New(&sb, DebugLevel, ZeroLogBackend)
	l.Info().AddTime(key, val).Flush("")
	s := sb.String()
	assert.Contains(t, s, key, "Message should contain key")
}
