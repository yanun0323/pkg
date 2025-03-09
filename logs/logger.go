package logs

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/yanun0323/pkg/logs/internal"
)

type loggerNew slog.Logger

func New(level Level, outputs ...io.Writer) Logger {
	var out io.Writer = os.Stdout
	if len(outputs) != 0 {
		out = outputs[0]
	}

	return (*loggerNew)(slog.New(internal.NewLoggerHandler(out, int8(level))))
}

func (l loggerNew) clone() *loggerNew {
	return (*loggerNew)((*slog.Logger)(&l))
}

func (l *loggerNew) Copy() Logger {
	return l.clone()
}

func (l *loggerNew) WithContext(ctx context.Context) Logger {
	return l.WithField("context", ctx)
}

func (l *loggerNew) withField(key string, value any) *loggerNew {
	return (*loggerNew)((*slog.Logger)(l).With(slog.Any(key, value)))
}

func (l *loggerNew) WithField(key string, value any) Logger {
	return l.withField(key, value)
}

func (l *loggerNew) WithFields(fields map[string]any) Logger {
	if len(fields) == 0 {
		return l
	}

	for k, v := range fields {
		l = l.withField(k, v)
	}

	return l
}

func (l *loggerNew) WithError(err error) Logger {
	return l.WithField("error", err)
}

func (l *loggerNew) Attach(ctx context.Context) (context.Context, Logger) {
	ll := l.clone()
	return context.WithValue(ctx, logKey{}, ll), ll
}

func (l *loggerNew) Log(level Level, args ...any) {
	(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprint(args...))
}

func (l *loggerNew) Logf(level Level, format string, args ...any) {
	if len(args) == 0 {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), format)
	} else {
		(*slog.Logger)(l).Log(context.Background(), slog.Level(level), fmt.Sprintf(format, args...))
	}
}

func (l *loggerNew) Debug(args ...any) {
	l.Log(LevelDebug, args...)
}

func (l *loggerNew) Debugf(format string, args ...any) {
	l.Logf(LevelDebug, format, args...)
}

func (l *loggerNew) Info(args ...any) {
	l.Log(LevelInfo, args...)
}

func (l *loggerNew) Infof(format string, args ...any) {
	l.Logf(LevelInfo, format, args...)
}

func (l *loggerNew) Warn(args ...any) {
	l.Log(LevelWarn, args...)
}

func (l *loggerNew) Warnf(format string, args ...any) {
	l.Logf(LevelWarn, format, args...)
}

func (l *loggerNew) Error(args ...any) {
	l.Log(LevelError, args...)
}

func (l *loggerNew) Errorf(format string, args ...any) {
	l.Logf(LevelError, format, args...)
}

func (l *loggerNew) Fatal(args ...any) {
	l.Log(LevelFatal, args...)
	os.Exit(1)
}

func (l *loggerNew) Fatalf(format string, args ...any) {
	l.Logf(LevelFatal, format, args...)
	os.Exit(1)
}
