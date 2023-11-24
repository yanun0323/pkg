package logs

import (
	"strings"

	"github.com/sirupsen/logrus"
)

type Level string

// Convert the Level to a string. E.g. PanicLevel becomes "panic".
func (level Level) String() string {
	return string(level)
}

// NewLevel takes a string level and returns the Logs log level constant.
//
// return panic level when there's no matched string
//
// allowed args: "panic", "fatal", "error", "warn", "info", "debug", "trace"
func NewLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "panic":
		return LevelPanic
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
	case "trace":
		return LevelTrace
	}

	return "panic"
}

func (level Level) logrus() logrus.Level {
	switch level {
	case LevelTrace:
		return logrus.TraceLevel
	case LevelDebug:
		return logrus.DebugLevel
	case LevelInfo:
		return logrus.InfoLevel
	case LevelWarn:
		return logrus.WarnLevel
	case LevelError:
		return logrus.ErrorLevel
	case LevelFatal:
		return logrus.FatalLevel
	case LevelPanic:
		return logrus.PanicLevel
	}

	return logrus.PanicLevel
}

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `logrus.New()`.
const (
	// LevelPanic level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic Level = "panic"
	// LevelFatal level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LevelFatal Level = "fatal"
	// LevelError level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError Level = "error"
	// LevelWarn level. Non-critical entries that deserve eyes.
	LevelWarn Level = "warn"
	// LevelInfo level. General operational entries about what's going on inside the
	// application.
	LevelInfo Level = "info"
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug Level = "debug"
	// LevelTrace level. Designates finer-grained informational events than the Debug.
	LevelTrace Level = "trace"
)
