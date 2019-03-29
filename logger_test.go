package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const IncorrectLevel Level = 1000

func tests() []struct {
	lvl Level
	key string
	val interface{}
} {
	return []struct {
		lvl Level
		key string
		val interface{}
	}{
		{DebugLevel, "lvldebug", "true"},
		{InfoLevel, "lvlinfo", "true"},
		{WarnLevel, "lvlwarn", "true"},
		{ErrorLevel, "lvlerr", "true"},
		{PanicLevel, "lvlpanic", "true"},
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		lvl    Level
		panics bool
	}{
		{DebugLevel, false},
		{InfoLevel, false},
		{WarnLevel, false},
		{ErrorLevel, false},
		{FatalLevel, false},
		{PanicLevel, false},
		{IncorrectLevel, true},
	}
	for _, test := range tests {
		var l Logger
		f1 := func() { l = New(os.Stdout, test.lvl, LogrusBackend) }
		f2 := func() { l = New(os.Stdout, test.lvl, ZeroLogBackend) }
		f3 := func() { l = New(os.Stdout, test.lvl, GelfBackend) }

		if test.panics {
			assert.Panics(t, f1, "Creating a logger with invalid level should panic")
			assert.Panics(t, f2, "Creating a logger with invalid level should panic")
			assert.Panics(t, f3, "Creating a logger with invalid level should panic")
		} else {
			assert.NotPanics(t, f1, "Creating a logger with valid level should work")
			assert.NotNil(t, l, "Logger should not be nil after creation")
			assert.NotPanics(t, f2, "Creating a logger with valid level should work")
			assert.NotNil(t, l, "Logger should not be nil after creation")
			assert.NotPanics(t, f3, "Creating a logger with valid level should work")
			assert.NotNil(t, l, "Logger should not be nil after creation")
		}
	}
}
