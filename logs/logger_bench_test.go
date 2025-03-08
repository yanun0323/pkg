package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	writerL   = FileWriter("./bench", "logger")
	writerLNH = FileWriter("./bench", "logger_new_handler")
	writerS   = FileWriter("./bench", "slog")
	writerZ   = FileWriter("./bench", "zap")
)

func BenchmarkLogger(b *testing.B) {
	New(LevelInfo, os.Stdout).WithField("key", "value").Info("start test")

	l := New(LevelInfo, writerL)
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
	slog.New(newLoggerHandler(os.Stdout, LevelInfo)).WithGroup("hello").With("key", "value").Error("start test")

	l := slog.New(newLoggerHandler(writerLNH, LevelInfo))
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

type colorHandler struct {
	out  io.Writer
	opts *slog.HandlerOptions
}

func NewColorHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	return &colorHandler{
		out:  w,
		opts: opts,
	}
}

func (h *colorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if h.opts == nil || h.opts.Level == nil {
		return true
	}

	return level >= h.opts.Level.Level()
}

func (h *colorHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := new(bytes.Buffer)

	// 時間格式
	timeStr := r.Time.Format("2006/01/02 15:04:05")
	buf.WriteString(timeStr)
	buf.WriteString(" ")

	// 彩色等級
	// 依據不同等級設定顏色
	// 拷貝記錄並修改 level 文字
	level := r.Level.String()

	// 依據不同等級設定顏色
	levelColor := colorize(level, NewLevel(level).color())
	buf.WriteString(levelColor)
	buf.WriteString(" ")

	// 訊息
	buf.WriteString(r.Message)
	buf.WriteString("  ")

	// 屬性
	r.Attrs(func(a slog.Attr) bool {
		// 跳過內部屬性
		if a.Key == "time" || a.Key == "level" || a.Key == "msg" {
			return true
		}
		buf.WriteString("[")
		buf.WriteString(a.Key)
		buf.WriteString("] ")
		buf.WriteString(fmt.Sprint(a.Value.Any()))
		buf.WriteString("  ")
		return true
	})

	// 上下文資訊
	buf.WriteString("[context] ")
	if ctx == nil {
		buf.WriteString("nil")
	} else {
		buf.WriteString("context.Background")
	}

	buf.WriteString("\n")

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *colorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// 簡化實作，實際使用中可能需要保存這些屬性
	return h
}

func (h *colorHandler) WithGroup(name string) slog.Handler {
	// 簡化實作，實際使用中可能需要處理群組
	return h
}
