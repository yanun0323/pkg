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

// 用於存儲群組或屬性的結構體
type groupOrAttrs struct {
	group string      // 群組名稱（如果非空）
	attrs []slog.Attr // 屬性（如果非空）
}

type loggerHandlerOptions struct {
	level Level
}

type loggerHandler struct {
	opts *loggerHandlerOptions
	goas []groupOrAttrs // 存儲 WithGroup 和 WithAttrs 的狀態
	out  io.Writer
}

func newLoggerHandler(w io.Writer, level Level) *loggerHandler {
	return &loggerHandler{
		opts: &loggerHandlerOptions{level: level},
		goas: []groupOrAttrs{},
		out:  w,
	}
}

func (h *loggerHandler) Level() slog.Level {
	return slog.Level(h.opts.level)
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

	// Process WithGroup and WithAttrs state
	goas := h.goas
	if r.NumAttrs() == 0 {
		// If the record has no attributes, remove the last empty group
		for len(goas) > 0 && goas[len(goas)-1].group != "" {
			goas = goas[:len(goas)-1]
		}
	}

	// Process groups and attributes
	for _, goa := range goas {
		if goa.group != "" {
			buf.WriteString(colorize("["+goa.group+"] ", colorMagenta))
		} else {
			for _, a := range goa.attrs {
				if len(a.Key) == 0 {
					continue
				}

				buf.WriteString(colorize("["+a.Key+"] ", colorBlue))
				buf.WriteString(colorize(fmt.Sprint(a.Value.Any()), colorBlack))
				buf.WriteString("  ")
			}
		}
	}

	// 屬性
	r.Attrs(func(a slog.Attr) bool {
		// 跳過內部屬性
		switch a.Key {
		case "", "time", "level", "msg", cloneKey:
			return true
		default:
			buf.WriteString(colorize("["+a.Key+"] ", colorMagenta))
			buf.WriteString(colorize(fmt.Sprint(a.Value.Any()), colorBlack))
			buf.WriteString("  ")
		}
		return true
	})

	buf.WriteString("\n")

	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *loggerHandler) clone() *loggerHandler {
	h2 := &loggerHandler{
		opts: h.opts,
		out:  h.out,
		goas: make([]groupOrAttrs, len(h.goas)+1),
	}
	copy(h2.goas, h.goas)
	return h2
}

// 創建一個新的 handler 並添加一個 groupOrAttrs
func (h *loggerHandler) withGroupOrAttrs(goa groupOrAttrs) *loggerHandler {
	h2 := h.clone()
	h2.goas[len(h2.goas)-1] = goa
	return h2
}

func (h *loggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{attrs: attrs})
}

func (h *loggerHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return h.withGroupOrAttrs(groupOrAttrs{group: name})
}
