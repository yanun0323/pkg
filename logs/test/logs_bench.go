package test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/yanun0323/pkg/logs"
	"github.com/yanun0323/pkg/logs/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func BenchmarkLogger(b *testing.B) {
	writerL := logs.FileWriter(".", "logger")
	logs.New(logs.LevelInfo, os.Stdout).WithField("key", "value").Info("start test")

	l := logs.New(logs.LevelInfo, writerL)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.WithField("key", "value").Info("test")
		}
	})

	b.Cleanup(func() {
		if err := writerL.Remove(); err != nil {
			b.Fatalf("remove writerL failed: %v", err)
		}
	})
}

func BenchmarkSlogWithNewHandler(b *testing.B) {
	writerLNH := logs.FileWriter(".", "logger_new_handler")
	slog.New(internal.NewLoggerHandler(os.Stdout, int8(logs.LevelInfo))).WithGroup("hello").With("key", "value").Error("start test")

	l := slog.New(internal.NewLoggerHandler(writerLNH, int8(logs.LevelInfo)))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.With("key", "value").Info("test")
		}
	})

	b.Cleanup(func() {
		if err := writerLNH.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

func BenchmarkSlog(b *testing.B) {
	writerS := logs.FileWriter(".", "slog")
	slog.New(slog.NewJSONHandler(os.Stdout, nil)).WithGroup("hello").With("key", "value").Error("start test")

	l := slog.New(slog.NewJSONHandler(writerS, nil))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.With("key", "value").Info("test")
		}
	})

	b.Cleanup(func() {
		if err := writerS.Remove(); err != nil {
			b.Fatalf("remove writerS failed: %v", err)
		}
	})
}

func BenchmarkZap(b *testing.B) {
	writerZ := logs.FileWriter(".", "zap")
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		os.Stdout,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	)).With(zap.Any("key", "value")).Info("start test")

	l := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(conf),
		writerZ,
		zap.NewAtomicLevelAt(zap.InfoLevel),
	))
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			l.With(zap.Any("key", "value")).Info("test")
		}
	})

	b.Cleanup(func() {
		if err := writerZ.Remove(); err != nil {
			b.Fatalf("remove writerZ failed: %v", err)
		}
	})
}
