package logs

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

type timeLog struct {
	logger   Logger
	interval time.Duration
	previous *atomic.Int64
}

func NewTimableLog(logger Logger, interval time.Duration) Logger {
	return &timeLog{
		logger:   logger.Copy(),
		interval: interval,
		previous: &atomic.Int64{},
	}
}

func (tl *timeLog) timeHook(fn func()) {
	now := time.Now()
	if now.Add(-tl.interval).UnixNano() < tl.previous.Load() {
		return
	}

	tl.previous.Store(now.UnixNano())
	fn()
}

func (tl *timeLog) SetOutput(output io.Writer) {
	tl.logger.SetOutput(output)
}

func (tl *timeLog) Debug(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Debug(args...)
	})
}

func (tl *timeLog) Debugf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Debugf(format, args...)
	})
}

func (tl *timeLog) Copy() Logger {
	p := atomic.Int64{}
	p.Store(tl.previous.Load())

	return &timeLog{
		logger:   tl.logger.Copy(),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) copyPrevious() *atomic.Int64 {
	a := &atomic.Int64{}
	a.Store(tl.previous.Load())
	return a
}

func (tl *timeLog) Error(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Error(args...)
	})
}

func (tl *timeLog) Errorf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Errorf(format, args...)
	})
}

func (tl *timeLog) Fatal(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Fatal(args...)
	})
}

func (tl *timeLog) Fatalf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Fatalf(format, args...)
	})
}

func (tl *timeLog) Info(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Info(args...)
	})
}

func (tl *timeLog) Infof(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Infof(format, args...)
	})
}

func (tl *timeLog) Log(level Level, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Log(level, args...)
	})
}

func (tl *timeLog) Logf(level Level, format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Logf(level, format, args...)
	})
}

func (tl *timeLog) Panic(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Panic(args...)
	})
}

func (tl *timeLog) Panicf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Panicf(format, args...)
	})
}

func (tl *timeLog) Print(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Print(args...)
	})
}

func (tl *timeLog) Printf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Printf(format, args...)
	})
}

func (tl *timeLog) Trace(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Trace(args...)
	})
}

func (tl *timeLog) Tracef(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Tracef(format, args...)
	})
}

func (tl *timeLog) Warn(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Warn(args...)
	})
}

func (tl *timeLog) Warnf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Warnf(format, args...)
	})
}

func (tl *timeLog) Warning(args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Warning(args...)
	})
}

func (tl *timeLog) Warningf(format string, args ...interface{}) {
	tl.timeHook(func() {
		tl.logger.Warningf(format, args...)
	})
}

func (tl *timeLog) WithContext(ctx context.Context) Logger {
	return &timeLog{
		logger:   tl.logger.WithContext(ctx),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) WithError(err error) Logger {
	return &timeLog{
		logger:   tl.logger.WithError(err),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) WithField(key string, value interface{}) Logger {
	return &timeLog{
		logger:   tl.logger.WithField(key, value),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) WithFields(fields map[string]interface{}) Logger {
	return &timeLog{
		logger:   tl.logger.WithFields(fields),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) WithTime(t time.Time) Logger {
	return &timeLog{
		logger:   tl.logger.WithTime(t),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) Writer() *io.PipeWriter {
	return tl.logger.Writer()
}

func (tl *timeLog) WriterLevel(level Level) *io.PipeWriter {
	return tl.logger.WriterLevel(level)
}

func (tl *timeLog) WithFunc(function string) Logger {
	return &timeLog{
		logger:   tl.logger.WithFunc(function),
		interval: tl.interval,
		previous: tl.copyPrevious(),
	}
}

func (tl *timeLog) Attach(ctx context.Context) (context.Context, Logger) {
	ll := tl.Copy()
	return context.WithValue(ctx, logKey{}, ll), ll
}
