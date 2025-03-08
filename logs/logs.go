package logs

import (
	"context"
	"sync/atomic"
)

var (
	// defaultLogger is the default logger.
	defaultLogger atomic.Value

	defaultTimeFormat atomic.Value
)

const (
	_defaultFormat = "2006/01/02 15:04:05"
)

func init() {
	defaultLogger.Store(New(LevelInfo))
	defaultTimeFormat.Store(_defaultFormat)
}

// logKey is the key for the logger in the context.
type logKey struct{}

// Get gets the logger from context. if there's no logger in context, it will create a new logger with 'info' level.
func Get(ctx context.Context) Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(Logger); ok {
		return logger
	}

	return Default()
}

// Default returns the default logger.
func Default() Logger {
	l, ok := defaultLogger.Load().(Logger)
	if !ok {
		l = New(LevelInfo)
		defaultLogger.Store(l)
	}

	return l
}

// SetDefault sets the default logger.
func SetDefault(logger Logger) {
	if logger != nil {
		defaultLogger.Store(logger)
	}
}

// SetDefaultTimeFormat sets the default time format.
func SetDefaultTimeFormat(format string) {
	if len(format) != 0 {
		defaultTimeFormat.Store(format)
	}
}

// Debug uses the default logger to log a message at the debug level.
func Debug(args ...interface{}) {
	Default().Debug(args...)
}

// Debugf uses the default logger to log a message at the debug level.
func Debugf(format string, args ...interface{}) {
	Default().Debugf(format, args...)
}

// Error uses the default logger to log a message at the error level.
func Error(args ...interface{}) {
	Default().Error(args...)
}

// Errorf uses the default logger to log a message at the error level.
func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

// Fatal uses the default logger to log a message at the fatal level.
func Fatal(args ...interface{}) {
	Default().Fatal(args...)
}

// Fatalf uses the default logger to log a message at the fatal level.
func Fatalf(format string, args ...interface{}) {
	Default().Fatalf(format, args...)
}

// Info uses the default logger to log a message at the info level.
func Info(args ...interface{}) {
	Default().Info(args...)
}

// Infof uses the default logger to log a message at the info level.
func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

// Warn uses the default logger to log a message at the warn level.
func Warn(args ...interface{}) {
	Default().Warn(args...)
}

// Warnf uses the default logger to log a message at the warn level.
func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}
