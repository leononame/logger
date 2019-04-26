# Logger [![GoDoc](https://godoc.org/github.com/leononame/logger?status.svg)](https://godoc.org/github.com/leononame/logger) [![Build Status](https://cloud.drone.io/api/badges/leononame/logger/status.svg)](https://cloud.drone.io/leononame/logger) [![codecov](https://codecov.io/gh/leononame/logger/branch/master/graph/badge.svg)](https://codecov.io/gh/leononame/logger)

This package contains an abstraction for logging so different log implementations can be used for different purposes.

## Why?

Depending on your use case you might want to use a different log implementation. Logrus looks nice on a console output, but is a rather slow JSON logger.   

Another reason might be you want to use multiple parametrized loggers. This package allows you to wrap multiple loggers around a single interface. Each logger can have its own io.Writer as Output and its own log level set. This way, you can send only Error+Above levels to one of your log servers and have a second log server which gets log levels of all type, plus a third logger which prints errors to stdout.

## Implementation

It's relatively easy to build a wrapper around the standard interface. There are currently three implementations: logrus, zerolog and gelf.

- [Logrus](https://github.com/leononame/logger-logrus) can be used as a console logger while running your application locally
- [Zerolog](https://github.com/leononame/logger-zerolog)  is a fast JSON logger that prints everything as a JSON message
- [GELF](https://github.com/leononame/logger-gelf)  is a wrapper around zerolog that prints everything in GELF format.

Additionally, a Multilogger exists. You can pass as many loggers as you want to your Multilogger and the Multilogger behaves as a single logger that calls the same functions on all Loggers.

## Usage

### Simple

Create a logger with new and select a predefined implementation.

```go
package main

import (
	"os"

	"github.com/juju/errors"
	"github.com/leononame/logger"
	"github.com/leononame/logger-logrus"
)

func main() {
	// Options are: LogrusBackend, ZerologBackend, GelfBackend
	l := logrus.New(os.Stdout, logger.InfoLevel)
	err := errors.New("test")
	l.Info().AddStr("key", "value").AddErr(err).Flush("some message")
}
```

Output:

```
INFO[0000] some message                                  err=test err_stack="/Users/leo/Documents/code/vgo/test/main.go:13: test" fields.time="2019-03-29 15:53:48.947011 -0500 -05 m=+0.000544729" key=value
```

### Multiple loggers

Use multiple loggers with different levels as if it was only a single one.

```go
package main

import (
	"os"

	"github.com/juju/errors"
	"github.com/leononame/logger"
	"github.com/leononame/logger-logrus"
	"github.com/leononame/logger-zerolog"
)

func main() {
	l1 := logrus.New(os.Stdout, logger.InfoLevel)
	l2 := zerolog.New(os.Stderr, logger.ErrorLevel)
	l := logger.NewMulti(l1, l2)
	// This gets printed to stdout with logrus
	l.Info().AddStr("key", "value").Flush("message")
	err := errors.New("err1")
	err = errors.Annotate(err, "additional trace")
	// This gets printed to stdout with logrus and stderr with zerolog
	l.Error().AddErr(err).Flush("additional message")
}
```

Output:

```
INFO[0000] message                                       fields.time="2019-03-29 15:57:51.307562 -0500 -05 m=+0.000549599" key=value
ERRO[0000] additional message                            err="additional trace: err1" err_stack="/Users/leo/Documents/code/vgo/test/main.go:16: err1\n/Users/leo/Documents/code/vgo/test/main.go:17: additional trace" fields.time="2019-03-29 15:57:51.307811 -0500 -05 m=+0.000799046"
{"level":"error","err":"additional trace: err1","err_stack":"/Users/leo/Documents/code/vgo/test/main.go:16: err1\n/Users/leo/Documents/code/vgo/test/main.go:17: additional trace","time":1553893071,"message":"additional message"}
```


### Customized logger

Customize your logrus parameters

```go
package main

import (
	"time"

	"github.com/leononame/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	ll := logrus.New()
	ll.WithField("source", "cli_client")
	ll.SetLevel(logrus.WarnLevel)
	// Alternatively: logger.FromZerolog instantiates a Logger with zerolog implementation
	l := logger.FromLogrus(ll)
	l.Info().AddStr("key", "value").Flush("message")
	l.Warn().AddInt("iteration", 1000).Flush("finished calculation")
	l.Error().AddDur("duration", time.Minute).Flush("duration calculated")
}
```

Output:

```
WARN[0000] finished calculation                          fields.time="2019-03-29 16:02:58.70876 -0500 -05 m=+0.000626253" iteration=1000
ERRO[0000] duration calculated                           duration=1m0s fields.time="2019-03-29 16:02:58.708996 -0500 -05 m=+0.000862119"
```

## Logger API

Each logger implements the interface below. Calling `WithField` returns a logger that always logs the specified field. All other calls return a log entry (not written yet) that will log at the specified level.

```go
// Logger is an standard interface for logging so that different log implementations can be wrapped around.
// The API is heavily influenced by the zerolog API for structured JSON logging
type Logger interface {
	// WithField returns a new Logger that always logs the specified field
	WithField(key, value string) Logger
	// Level creates a new Entry with the specified Level
	Level(Level) Entry
	// Debug creates a new Entry with level Debug
	Debug() Entry
	// Info creates a new Entry with level Info
	Info() Entry
	// Warn creates a new Entry with level Warn
	Warn() Entry
	// Error creates a new Entry with level Error
	Error() Entry
	// Fatal creates a new Entry with level Fatal. Executing a log at fatal level exits the application with exit code 1.
	Fatal() Entry
	// Panic creates a new Entry with level Panic. Executing a log at panic level will call panic().
	Panic() Entry
}
```

Below you will find the interface of a log entry. There are different functinos to add custom fields for structured logging. The `Flush` function dumps the log (only if the level of the entry is at least the logger's level). If `Flush` is not called, nothing will be logged

```go
// Entry is an interface for a log entry. A single entry always has defined a log level. Custom fields can be
// added. An entry will never be written to the log unless Flush is called.
type Entry interface {
	// Flush writes the entry as a single log statement. Optionally, a message can be added which will
	// be included in the final log entry
	Flush(string)

	// Add a range of fields to the log statement
	AddFields(map[string]interface{}) Entry
	// Add an error to the log statement. The error will have the key "err". An error stack will be included
	// under the key "err_stack"
	AddErr(err error) Entry
	// Add an error to the log statement. An error stack will be included under the key "${key}_stack"
	AddError(key string, val error) Entry
	// Add a bool value to the log statement.
	AddBool(key string, val bool) Entry
	// Add an integer value to the log statement.
	AddInt(key string, val int) Entry
	// Add a string value to the log statement.
	AddStr(key string, val string) Entry
	// Add a time value to the log statement.
	AddTime(key string, val time.Time) Entry
	// Add a duration value to the log statement.
	AddDur(key string, val time.Duration) Entry
	// Add any value to the log statement.
	AddAny(key string, val interface{}) Entry
}
```