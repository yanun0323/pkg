package logs

import (
	"log/slog"
	"testing"
)

type EmptyWriter struct{}

func (EmptyWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// go test -bench=. -benchmem ./...
func BenchmarkLogger(b *testing.B) {
	l := New(LevelInfo, EmptyWriter{})
	for i := 0; i < b.N; i++ {
		l.WithField("key", "value").Info("test")
	}
}

func BenchmarkLoggerLight(b *testing.B) {
	l := NewLightweightLogger(LevelInfo, EmptyWriter{})
	for i := 0; i < b.N; i++ {
		l.WithField("key", "value").Info("test")
	}
}

func BenchmarkSlog(b *testing.B) {
	l := slog.New(slog.NewJSONHandler(EmptyWriter{}, nil))
	for i := 0; i < b.N; i++ {
		l.With(slog.Any("key", "value")).Info("test")
	}
}
