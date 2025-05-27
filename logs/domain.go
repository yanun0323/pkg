// Package logs is deprecated and no longer maintained.
// Please use github.com/yanun0323/logs as an alternative.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
package logs

import (
	"context"
)

// Logger is the interface that wraps the basic methods of a logger.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
type Logger interface {
	// Copy duplicates the logger.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Copy() Logger

	// Attach attaches the logger into the context.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Attach(ctx context.Context) (context.Context, Logger)

	// WithField adds a single field to the Logger.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	WithField(key string, value interface{}) Logger
	// WithFields adds a map of fields to the Logger.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	WithFields(fields map[string]interface{}) Logger
	// WithError adds an error as single field (using the key defined in ErrorKey) to the Logger.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	WithError(err error) Logger
	// WithContext adds a context to the Logger.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	WithContext(ctx context.Context) Logger

	// Log will log a message at the level given as parameter.
	// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
	// For this behavior Entry.Panic or Entry.Fatal should be used instead.	Log(level Level, args ...interface{})
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Log(level Level, args ...interface{})

	// Logf will log a message at the level given as parameter.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Logf(level Level, format string, args ...interface{})

	// Debug will log a message at the debug level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Debug(args ...interface{})

	// Debugf will log a message at the debug level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Debugf(format string, args ...interface{})

	// Info will log a message at the info level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Info(args ...interface{})

	// Infof will log a message at the info level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Infof(format string, args ...interface{})

	// Warn will log a message at the warn level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Warn(args ...interface{})

	// Warnf will log a message at the warn level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Warnf(format string, args ...interface{})

	// Error will log a message at the error level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Error(args ...interface{})

	// Errorf will log a message at the error level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Errorf(format string, args ...interface{})

	// Fatal will log a message at the fatal level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Fatal(args ...interface{})

	// Fatalf will log a message at the fatal level.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	Fatalf(format string, args ...interface{})
}
