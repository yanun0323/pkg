package logs

import (
	"context"
	"io"
	"sync"
	"time"
)

type timerLoggerNew struct {
	mu                  sync.Mutex
	last                int64
	intervalMillisecond int64

	Logger
}

// NewTimerLogger creates a new timer logger with the given interval and level.
//
// A timer logger is a logger that logs messages only when the interval has passed.
//
// If outputs is not provided, the logger will write to the os.Stdout.
func NewTimerLogger(interval time.Duration, level Level, outputs ...io.Writer) Logger {
	itv := interval.Milliseconds()
	return &timerLoggerNew{
		last:                time.Now().UnixMilli() - itv,
		intervalMillisecond: itv,
		Logger:              New(level, outputs...),
	}
}

func (l *timerLoggerNew) canBeFire() bool {
	now := time.Now().UnixMilli()
	if !l.mu.TryLock() {
		return false
	}

	available := l.last + l.intervalMillisecond
	ok := now >= available
	if ok {
		l.last = now
	}
	l.mu.Unlock()

	return ok
}

func (l *timerLoggerNew) Copy() Logger {
	copied := &timerLoggerNew{
		last:                l.last,
		intervalMillisecond: l.intervalMillisecond,
		Logger:              l.Logger.Copy(),
	}
	return copied
}

func (l *timerLoggerNew) WithContext(ctx context.Context) Logger {
	return l.Logger.WithContext(ctx)
}

func (l *timerLoggerNew) WithField(key string, value any) Logger {
	return l.Logger.WithField(key, value)
}

func (l *timerLoggerNew) WithFields(fields map[string]any) Logger {
	return l.Logger.WithFields(fields)
}

func (l *timerLoggerNew) WithError(err error) Logger {
	return l.Logger.WithError(err)
}

func (l *timerLoggerNew) Attach(ctx context.Context) (context.Context, Logger) {
	return l.Logger.Attach(ctx)
}

func (l *timerLoggerNew) Log(level Level, args ...any) {
	if l.canBeFire() {
		l.Logger.Log(level, args...)
	}
}

func (l *timerLoggerNew) Logf(level Level, format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Logf(level, format, args...)
	}
}

func (l *timerLoggerNew) Debug(args ...any) {
	if l.canBeFire() {
		l.Logger.Debug(args...)
	}
}

func (l *timerLoggerNew) Debugf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Debugf(format, args...)
	}
}

func (l *timerLoggerNew) Info(args ...any) {
	if l.canBeFire() {
		l.Logger.Info(args...)
	}
}

func (l *timerLoggerNew) Infof(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Infof(format, args...)
	}
}

func (l *timerLoggerNew) Warn(args ...any) {
	if l.canBeFire() {
		l.Logger.Warn(args...)
	}
}

func (l *timerLoggerNew) Warnf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Warnf(format, args...)
	}
}

func (l *timerLoggerNew) Error(args ...any) {
	if l.canBeFire() {
		l.Logger.Error(args...)
	}
}

func (l *timerLoggerNew) Errorf(format string, args ...any) {
	if l.canBeFire() {
		l.Logger.Errorf(format, args...)
	}
}

func (l *timerLoggerNew) Fatal(args ...any) {
	l.Logger.Fatal(args...)
}

func (l *timerLoggerNew) Fatalf(format string, args ...any) {
	l.Logger.Fatalf(format, args...)
}
