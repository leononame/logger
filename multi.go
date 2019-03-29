package logger

import (
	"time"
)

func NewMulti(ls ...Logger) Logger {
	return &mLog{ls: ls}
}

type mLog struct {
	ls []Logger
}

// WithField returns a new Logger that always logs the specified field
func (m *mLog) WithField(key, value string) Logger {
	for i := range m.ls {
		m.ls[i] = m.ls[i].WithField(key, value)
	}
	return m
}

// Level creates a new Entry with the specified Level
func (m *mLog) Level(lvl Level) Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Level(lvl)
	}
	return &e
}

// Debug creates a new Entry with level Debug
func (m *mLog) Debug() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Debug()
	}
	return &e
}

// Info creates a new Entry with level Info
func (m *mLog) Info() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Info()
	}
	return &e
}

// Warn creates a new Entry with level Warn
func (m *mLog) Warn() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Warn()
	}
	return &e
}

// Error creates a new Entry with level Error
func (m *mLog) Error() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Error()
	}
	return &e
}

// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
func (m *mLog) Fatal() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Fatal()
	}
	return &e
}

// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
func (m *mLog) Panic() Entry {
	e := mEntry{make([]Entry, len(m.ls))}
	for i := range m.ls {
		e.es[i] = m.ls[i].Panic()
	}
	return &e
}

type mEntry struct {
	es []Entry
}

// Flush writes the entry as a single log statement. Optionally, a message can be added which will
// be included in the final log entry
func (m *mEntry) Flush(msg string) {
	var r interface{} = nil
	// Since this might be a Flush() call on panic level, we have to recover from each call.
	// Otherwise, the rest of the loggers will not flush
	for i := range m.es {
		func() {
			defer func() {
				r = recover()
			}()
			m.es[i].Flush(msg)
		}()
	}
	// If any panic call was recovered from, panic again, but only once
	if r != nil {
		panic(r)
	}
}

// AddFields adds a range of fields to the log statement
func (m *mEntry) AddFields(fields map[string]interface{}) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddFields(fields)
	}
	return m
}

// AddErr adds an error to the log statement. The error will have the key "err". An error stack will be included
// under the key "err_stack"
func (m *mEntry) AddErr(err error) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddErr(err)
	}
	return m
}

// AddError adds an error to the log statement. An error stack will be included under the key "${key}_stack"
func (m *mEntry) AddError(key string, val error) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddError(key, val)
	}
	return m
}

// AddBool adds a bool value to the log statement.
func (m *mEntry) AddBool(key string, val bool) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddBool(key, val)
	}
	return m
}

// AddInt adds an integer value to the log statement.
func (m *mEntry) AddInt(key string, val int) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddInt(key, val)
	}
	return m
}

// AddStr adds a string value to the log statement.
func (m *mEntry) AddStr(key string, val string) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddStr(key, val)
	}
	return m
}

// AddTime adds a time value to the log statement.
func (m *mEntry) AddTime(key string, val time.Time) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddTime(key, val)
	}
	return m
}

// AddDur adds a duration value to the log statement.
func (m *mEntry) AddDur(key string, val time.Duration) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddDur(key, val)
	}
	return m
}

// AddAny adds any value to the log statement.
func (m *mEntry) AddAny(key string, val interface{}) Entry {
	for i := range m.es {
		m.es[i] = m.es[i].AddAny(key, val)
	}
	return m
}
