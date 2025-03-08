package logs

import (
	"context"
)

type Logger interface {
	// Copy duplicates the logger.
	Copy() Logger

	// Attach attaches the logger into the context.
	Attach(ctx context.Context) (context.Context, Logger)

	// WithField adds a single field to the Logger.
	WithField(key string, value interface{}) Logger
	// WithFields adds a map of fields to the Logger.
	WithFields(fields map[string]interface{}) Logger
	// WithError adds an error as single field (using the key defined in ErrorKey) to the Logger.
	WithError(err error) Logger
	// WithContext adds a context to the Logger.
	WithContext(ctx context.Context) Logger

	// Log will log a message at the level given as parameter.
	// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
	// For this behavior Entry.Panic or Entry.Fatal should be used instead.	Log(level Level, args ...interface{})
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}
