package internal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
)

var (
	cloneKey = "!LOGS_CLONE"
)

type loggerHandler struct {
	level *int8
	buf   bytes.Buffer
	out   io.Writer
}

func NewLoggerHandler(w io.Writer, level int8) *loggerHandler {
	return &loggerHandler{
		level: &level,
		out:   w,
	}
}

func (h *loggerHandler) Level() slog.Level {
	return slog.Level(*h.level)
}

func (h *loggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.Level()
}

func (h *loggerHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := new(bytes.Buffer)
	// Time format
	buf.WriteString(Colorize(r.Time.Format(GetDefaultTimeFormat()), colorBlack))
	buf.WriteString(" ")

	// Colorize level
	level := int8(r.Level)
	buf.WriteString(Colorize(LevelTitle(level), LevelColor(level)))
	buf.WriteString(" ")

	// Message
	buf.WriteString(r.Message)
	buf.WriteString("  ")

	buf.Write(h.buf.Bytes())

	buf.WriteString("\n")

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h loggerHandler) clone() *loggerHandler {
	return &h
}

func (h *loggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	hh := h.clone()
	for _, a := range attrs {
		if a.Key == cloneKey {
			continue
		}
		hh.buf.WriteString(Colorize("["+a.Key+"] ", colorMagenta))
		hh.buf.WriteString(Colorize(fmt.Sprint(a.Value.Any()), colorBlack))
		hh.buf.WriteString("  ")
	}

	return hh
}

func (h *loggerHandler) WithGroup(name string) slog.Handler {
	return h
}

const (
	LevelFatal int8 = 12
	LevelError int8 = 8
	LevelWarn  int8 = 4
	LevelInfo  int8 = 0
	LevelDebug int8 = -4
)

func LevelTitle(level int8) string {
	switch level {
	case LevelFatal:
		return "FATAL"
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	}

	return "INFO"
}

func LevelColor(level int8) int {
	switch level {
	case LevelDebug:
		return colorBlue
	case LevelInfo:
		return colorGreen
	case LevelError:
		return colorDarkRed
	case LevelWarn:
		return colorYellow
	case LevelFatal:
		return colorBrightRed
	default:
		return colorBlue
	}
}
