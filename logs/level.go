package logs

import (
	"log/slog"
	"strings"
)

type Level int8

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	switch level {
	case LevelFatal:
		return "fatal"
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebug:
		return "debug"
	}

	return "panic"
}

// NewLevel takes a string level and returns the Logs log level constant.
//
// return panic level when there's no matched string
//
// allowed args: "panic", "fatal", "error", "warn", "info", "debug", "trace"
func NewLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "fatal":
		return LevelFatal
	case "error":
		return LevelError
	case "warn", "warning":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	}

	return LevelFatal
}

const (
	// LevelInfo level. General operational entries about what's going on inside the
	// application.
	LevelInfo Level = Level(slog.LevelInfo)
	// LevelWarn level. Non-critical entries that deserve eyes.
	LevelWarn Level = Level(slog.LevelWarn)
	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError Level = Level(slog.LevelError)
	// LevelFatal level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LevelFatal Level = Level(12)
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug Level = Level(slog.LevelDebug)
)
