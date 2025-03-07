package logs

import "sync/atomic"

var _defaultLogger atomic.Value

func Default() Logger {
	l, ok := _defaultLogger.Load().(Logger)
	if !ok {
		l = New(LevelInfo)
		_defaultLogger.Store(l)
	}

	return l
}

func SetDefault(logger Logger) {
	_defaultLogger.Store(logger)
}

func SetDefaultLevel(level Level) {
	_defaultLogger.Store(New(level))
}

func Debug(args ...interface{}) {
	Default().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Default().Debugf(format, args...)
}

func Error(args ...interface{}) {
	Default().Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Default().Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Default().Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Default().Fatalf(format, args...)
}

func Info(args ...interface{}) {
	Default().Info(args...)
}

func Infof(format string, args ...interface{}) {
	Default().Infof(format, args...)
}

func Panic(args ...interface{}) {
	Default().Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	Default().Panicf(format, args...)
}

func Trace(args ...interface{}) {
	Default().Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	Default().Tracef(format, args...)
}

func Warn(args ...interface{}) {
	Default().Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Default().Warnf(format, args...)
}

func Warning(args ...interface{}) {
	Default().Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	Default().Warningf(format, args...)
}
