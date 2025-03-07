package logs

import (
	"log/slog"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type EmptyWriter struct{}

func (EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (EmptyWriter) Sync() error {
	return nil
}

// go test -bench=. -benchmem ./...
func BenchmarkLogger(b *testing.B) {
	l := New(LevelInfo, EmptyWriter{})
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.WithField("key", "value").Info("test")
		}
	})
}

func BenchmarkLoggerLight(b *testing.B) {
	l := NewLightweightLogger(LevelInfo, EmptyWriter{})
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.WithField("key", "value").Info("test")
		}
	})
}

func BenchmarkSlog(b *testing.B) {
	l := slog.New(slog.NewJSONHandler(EmptyWriter{}, nil))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.With("key", "value").Info("test")
		}
	})
}

func BenchmarkZap(b *testing.B) {
	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
		EmptyWriter{},
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.With(zap.Any("key", "value")).Info("test")
		}
	})
}
