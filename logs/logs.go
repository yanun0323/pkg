package logs

import (
	"context"

	"github.com/yanun0323/pkg/logs/internal"
)

var (
	// defaultLogger is the default logger.
	defaultLogger = internal.NewValue(New(LevelInfo))
)

// logKey is the key for the logger in the context.
type logKey struct{}

// Get gets the logger from context. if there's no logger in context, it will create a new logger with 'info' level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Get(ctx context.Context) Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(Logger); ok {
		return logger
	}

	return Default()
}

// Default returns the default logger.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Default() Logger {
	l, ok := defaultLogger.Load().(Logger)
	if !ok {
		l = New(LevelInfo)
		defaultLogger.Store(l)
	}

	return l
}

// SetDefault sets the default logger.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func SetDefault(logger Logger) {
	if logger != nil {
		defaultLogger.Store(logger)
	}
}

// SetDefaultTimeFormat sets the default time format.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func SetDefaultTimeFormat(format string) {
	if len(format) != 0 {
		internal.SetDefaultTimeFormat(format)
	}
}

// Debug uses the default logger to log a message at the debug level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Debug(args ...interface{}) {
	Default().Debug(args...)
}

// Debugf uses the default logger to log a message at the debug level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Debugf(format string, args ...interface{}) {
	Default().Debugf(format, args...)
}

// Error uses the default logger to log a message at the error level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Error(args ...interface{}) {
	Default().Error(args...)
}

// Errorf uses the default logger to log a message at the error level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

// Fatal uses the default logger to log a message at the fatal level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Fatal(args ...interface{}) {
	Default().Fatal(args...)
}

// Fatalf uses the default logger to log a message at the fatal level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Fatalf(format string, args ...interface{}) {
	Default().Fatalf(format, args...)
}

// Info uses the default logger to log a message at the info level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Info(args ...interface{}) {
	Default().Info(args...)
}

// Infof uses the default logger to log a message at the info level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

// Warn uses the default logger to log a message at the warn level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Warn(args ...interface{}) {
	Default().Warn(args...)
}

// Warnf uses the default logger to log a message at the warn level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}
