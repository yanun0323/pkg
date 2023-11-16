package logs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var _defaultLogger *logger

type logKey struct{}

/*
Get logger from context. if there's no logger in context, it will create one.
*/
func Get(ctx context.Context) Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(*logger); ok {
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
	  0 = "panic"
	  1 = "fatal"
	  2 = "error"
	  3 = "warn"
	  4 = "info"
	  5 = "debug"
	  6 = "trace"
*/
func New(service string, level uint16) Logger {
	return NewLogger(service, level, "stdout")
}

/*
Init a new logger for output.

	# level
	  0 = "panic"
	  1 = "fatal"
	  2 = "error"
	  3 = "warn"
	  4 = "info"
	  5 = "debug"
	  6 = "trace"
*/
func NewLogger(service string, level uint16, outs string) Logger {
	l := logrus.New()
	l.SetLevel(Level(level).Logrus())
	l.SetNoLock()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		DataKey:         "fields",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "@timestamp",
		},
	})

	var out output
	l.SetOutput(out.new(outs, service))

	if len(service) != 0 {
		l.WithField("service", service)
	}

	log := &logger{
		Entry: l.WithContext(context.Background()),
	}

	if _defaultLogger == nil {
		_defaultLogger = log
	}

	return log
}

type logger struct {
	*logrus.Entry
}

func (l *logger) Copy() Logger {
	copied := l.Entry.Dup()
	return &logger{
		Entry: copied,
	}
}

func (l *logger) Log(level Level, args ...interface{}) {
	l.Entry.Log(level.Logrus(), args...)
}

func (l *logger) Logf(level Level, format string, args ...interface{}) {
	l.Entry.Logf(level.Logrus(), format, args...)
}

func (l *logger) Logln(level Level, args ...interface{}) {
	l.Entry.Logln(level.Logrus(), args...)
}

func (l *logger) WithContext(ctx context.Context) Logger {
	return &logger{
		Entry: l.Entry.WithContext(ctx),
	}
}

func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l
	}

	return &logger{
		Entry: l.Entry.WithField("error", fmt.Sprintf("%+v", err)),
	}
}

func (l *logger) WithField(k string, v interface{}) Logger {
	return &logger{
		Entry: l.Entry.WithField(k, v),
	}
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{
		Entry: l.Entry.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) WithTime(t time.Time) Logger {
	return &logger{
		Entry: l.Entry.WithTime(t),
	}
}

func (l *logger) WriterLevel(level Level) *io.PipeWriter {
	return l.Entry.WriterLevel(level.Logrus())
}

func (l *logger) WithFunc(function string) *logger {
	return &logger{
		Entry: l.Entry.WithField("func", function),
	}
}

func (l *logger) Attach(ctx context.Context) (Logger, context.Context) {
	return &logger{
			Entry: l.Entry,
		}, context.WithValue(ctx, logKey{}, &logger{
			Entry: l.Entry,
		})
}
