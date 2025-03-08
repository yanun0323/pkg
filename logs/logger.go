package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type loggerNew slog.Logger

func New(level Level, outputs ...io.Writer) Logger {
	var out io.Writer = os.Stdout
	if len(outputs) != 0 {
		out = outputs[0]
	}

	return (*loggerNew)(slog.New(newLoggerHandler(out, level)))
}

func (l *loggerNew) clone() *loggerNew {
	return (*loggerNew)((*slog.Logger)(l).With(slog.Any(cloneKey, struct{}{})))
}

func (l *loggerNew) Copy() Logger {
	return l.clone()
}

func (l *loggerNew) WithContext(ctx context.Context) Logger {
	return l.WithField("context", ctx)
}

func (l *loggerNew) WithField(key string, value interface{}) Logger {
	return (*loggerNew)((*slog.Logger)(l).With(slog.Any(key, value)))
}

func (l *loggerNew) WithFields(fields map[string]interface{}) Logger {
	if len(fields) == 0 {
		return l
	}

	attrs := make([]any, 0, len(fields))
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}

	return (*loggerNew)((*slog.Logger)(l).With(attrs...))
}

func (l *loggerNew) WithError(err error) Logger {
	return l.WithField("error", err)
}

func (l *loggerNew) Attach(ctx context.Context) (context.Context, Logger) {
	ll := l.clone()
	return context.WithValue(ctx, logKey{}, ll), ll
}

func (l *loggerNew) Log(level Level, args ...interface{}) {
	(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprint(args...))
}

func (l *loggerNew) Logf(level Level, format string, args ...interface{}) {
	if len(args) == 0 {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), format)
	} else {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprintf(format, args...))
	}
}

func (l *loggerNew) Debug(args ...interface{}) {
	l.Log(LevelDebug, args...)
}

func (l *loggerNew) Debugf(format string, args ...interface{}) {
	l.Logf(LevelDebug, format, args...)
}

func (l *loggerNew) Info(args ...interface{}) {
	l.Log(LevelInfo, args...)
}

func (l *loggerNew) Infof(format string, args ...interface{}) {
	l.Logf(LevelInfo, format, args...)
}

func (l *loggerNew) Warn(args ...interface{}) {
	l.Log(LevelWarn, args...)
}

func (l *loggerNew) Warnf(format string, args ...interface{}) {
	l.Logf(LevelWarn, format, args...)
}

func (l *loggerNew) Error(args ...interface{}) {
	l.Log(LevelError, args...)
}

func (l *loggerNew) Errorf(format string, args ...interface{}) {
	l.Logf(LevelError, format, args...)
}

func (l *loggerNew) Fatal(args ...interface{}) {
	l.Log(LevelFatal, args...)
	os.Exit(1)
}

func (l *loggerNew) Fatalf(format string, args ...interface{}) {
	l.Logf(LevelFatal, format, args...)
	os.Exit(1)
}
