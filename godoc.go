// Package logger contains an interface for abstracting any log implementation in case the implementation
// needs to be switched.
//
// Logging can be a performance bottleneck due to slow JSON marshalling or bad concurrent implementation. Hence,
// an abstraction is needed. Currently this package implements two different log backends, zerolog for fast
// JSON logging and logrus for pretty logging. The implementation can be chosen on creation.
package logger
