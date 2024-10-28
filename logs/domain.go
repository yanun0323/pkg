package logs

import (
	"context"
	"io"
	"time"
)

type Logger interface {
	// SetOutput sets the logger output.
	SetOutput(output io.Writer)
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	// Copy duplicates the logger.
	Copy() Logger
	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	// Log will log a message at the level given as parameter.
	// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
	// For this behavior Entry.Panic or Entry.Fatal should be used instead.	Log(level Level, args ...interface{})
	Log(level Level, args ...interface{})
	Logf(level Level, format string, args ...interface{})

	Panic(args ...interface{})
	Panicf(format string, args ...interface{})

	Print(args ...interface{})
	Printf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warning(args ...interface{})
	Warningf(format string, args ...interface{})

	// WithContext adds a context to the Logger.
	WithContext(ctx context.Context) Logger
	// WithError adds an error as single field (using the key defined in ErrorKey) to the Logger.
	WithError(err error) Logger
	// WithField adds a single field to the Logger.
	WithField(key string, value interface{}) Logger
	// WithFields adds a map of fields to the Logger.
	WithFields(fields map[string]interface{}) Logger
	// WithTime overrides the time of the Logger.
	WithTime(t time.Time) Logger
	Writer() *io.PipeWriter
	WriterLevel(level Level) *io.PipeWriter

	// WithField adds a func field to the Logger.
	WithFunc(function string) Logger
	// Attach attaches the logger into the context.
	Attach(ctx context.Context) (context.Context, Logger)
}
