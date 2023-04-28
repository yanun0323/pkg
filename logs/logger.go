package logs

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

var _defaultLogger *Logger

type logKey struct{}

/*
Get logger from context. if there's no logger in context, it will create one.
*/
func Get(ctx context.Context) *Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(*Logger); ok {
		return logger
	}
	if _defaultLogger == nil {
		return NewLogger("default", 2, "stdout")
	}
	return _defaultLogger
}

/*
Init a new logger for output.

	# level
	  0 = "info"
	  1 = "trace"
	  2 = "debug"
	  3 = "warn"
	  4 = "error"
	  5 = "fatal"
*/
func New(app string, level uint8) *Logger {
	return NewLogger(app, level, "stdout")
}

/*
Init a new logger for output.

	# level
	  0 = "info"
	  1 = "trace"
	  2 = "debug"
	  3 = "warn"
	  4 = "error"
	  5 = "fatal"
*/
func NewLogger(app string, level uint8, outs string) *Logger {
	logger := logrus.New()

	switch level {
	case 1:
		logger.SetLevel(logrus.TraceLevel)
	case 2:
		logger.SetLevel(logrus.DebugLevel)
	case 3:
		logger.SetLevel(logrus.WarnLevel)
	case 4:
		logger.SetLevel(logrus.ErrorLevel)
	case 5:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.SetNoLock()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		DataKey:         "fields",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "datetime",
		},
	})

	var out output
	logger.SetOutput(out.new(outs, app))

	if len(app) != 0 {
		logger.WithField("app", app)
	}

	log := &Logger{
		Entry: logger.WithContext(nil),
	}
	if _defaultLogger == nil {
		_defaultLogger = log
	}
	return log
}

type Logger struct {
	*logrus.Entry
}

func (l *Logger) WithEventID(eventID interface{}) *Logger {
	return &Logger{
		Entry: l.Entry.WithField("eventId", eventID),
	}
}

func (l *Logger) WithSource(source int) *Logger {
	return &Logger{
		Entry: l.Entry.WithField("source", source),
	}
}

func (l *Logger) WithFunc(function string) *Logger {
	return &Logger{
		Entry: l.Entry.WithField("func", function),
	}
}

func (l *Logger) WithField(k string, v interface{}) *Logger {
	return &Logger{
		Entry: l.Entry.WithField(k, v),
	}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	return &Logger{
		Entry: l.Entry.WithFields(logrus.Fields(fields)),
	}
}

func (l *Logger) WithUserID(userID uint64) *Logger {
	return &Logger{
		Entry: l.Entry.WithField("userId", userID),
	}
}

func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}

	return &Logger{
		Entry: l.Entry.WithField("error", fmt.Sprintf("%+v", err)),
	}
}

func (l *Logger) Attach(ctx context.Context) (*Logger, context.Context) {
	return &Logger{
			Entry: l.Entry,
		}, context.WithValue(ctx, logKey{}, &Logger{
			Entry: l.Entry,
		})
}
