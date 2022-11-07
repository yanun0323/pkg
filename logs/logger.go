package logs

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

var _DefaultLogger *Logger

type logKey struct{}

/*
Get logger from context. if there's no logger in context, it will create one.
*/
func Get(ctx context.Context) *Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(*Logger); ok {
		return logger
	}
	if _DefaultLogger == nil {
		return NewLogger("default", 2, "stdout")
	}
	return _DefaultLogger
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
			logrus.FieldKeyTime: "datatime",
		},
	})

	var out output
	logger.SetOutput(out.new(outs, app))

	log := &Logger{
		Entry: logger.WithField("app", app),
	}
	if _DefaultLogger == nil {
		_DefaultLogger = log
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

func (l *Logger) WithSlackNotify() *Logger {
	return &Logger{
		Entry: l.Entry.WithField("slack.notify", true),
	}
}

func (l *Logger) WithFunc(function string) *Logger {
	return &Logger{
		Entry: l.Entry.WithField("func", function),
	}
}

func (l *Logger) WithPair(base, quote interface{}) *Logger {
	base = strings.ToUpper(base.(string))
	quote = strings.ToUpper(quote.(string))

	return &Logger{
		Entry: l.Entry.WithField("base", base).WithField("quote", quote).WithField("pair", fmt.Sprintf("%s_%s", base, quote)),
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

// func Cause(err error) error {
// 	type causer interface {
// 		Cause() error
// 	}

// 	cause, ok := err.(causer)
// 	if !ok {
// 		fmt.Println(err)
// 		return err
// 	}
// 	err = cause.Cause()
// 	return err
// }
