package logs

import (
	"log/slog"
	"strings"
)

// Level is the type of the log level.
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
type Level int8

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
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
//
// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
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
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	LevelInfo Level = Level(slog.LevelInfo)
	// LevelWarn level. Non-critical entries that deserve eyes.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	LevelWarn Level = Level(slog.LevelWarn)
	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	LevelError Level = Level(slog.LevelError)
	// LevelFatal level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	LevelFatal Level = Level(12)
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	//
	// Deprecated: This package has been discontinued. Use github.com/yanun0323/logs instead.
	LevelDebug Level = Level(slog.LevelDebug)
)
