package logs

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
	level *Level
	buf   bytes.Buffer
	out   io.Writer
}

func newLoggerHandler(w io.Writer, level Level) *loggerHandler {
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
	buf.WriteString(colorize(r.Time.Format(defaultTimeFormat.Load().(string)), colorBlack))
	buf.WriteString(" ")

	// Colorize level
	level := Level(r.Level)
	buf.WriteString(colorize(level.title(), level.color()))
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
		hh.buf.WriteString(colorize("["+a.Key+"] ", colorMagenta))
		hh.buf.WriteString(colorize(fmt.Sprint(a.Value.Any()), colorBlack))
		hh.buf.WriteString("  ")
	}

	return hh
}

func (h *loggerHandler) WithGroup(name string) slog.Handler {
	return h
}
