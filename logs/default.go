package logs

func Debug(args ...interface{}) {
	_defaultLogger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	_defaultLogger.Debugf(format, args...)
}

func Error(args ...interface{}) {
	_defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	_defaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	_defaultLogger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	_defaultLogger.Fatalf(format, args...)
}

func Info(args ...interface{}) {
	_defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	_defaultLogger.Infof(format, args...)
}

func Panic(args ...interface{}) {
	_defaultLogger.Panic(args...)
}

func Panicf(format string, args ...interface{}) {
	_defaultLogger.Panicf(format, args...)
}

func Print(args ...interface{}) {
	_defaultLogger.Print(args...)
}

func Printf(format string, args ...interface{}) {
	_defaultLogger.Printf(format, args...)
}

func Trace(args ...interface{}) {
	_defaultLogger.Trace(args...)
}

func Tracef(format string, args ...interface{}) {
	_defaultLogger.Tracef(format, args...)
}

func Warn(args ...interface{}) {
	_defaultLogger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	_defaultLogger.Warnf(format, args...)
}

func Warning(args ...interface{}) {
	_defaultLogger.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	_defaultLogger.Warningf(format, args...)
}
