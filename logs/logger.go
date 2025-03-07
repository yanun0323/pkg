package logs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

type logKey struct{}

// Get gets the logger from context. if there's no logger in context, it will create a new logger with 'info' level.
func Get(ctx context.Context) Logger {
	val := ctx.Value(logKey{})
	if logger, ok := val.(Logger); ok {
		return logger
	}

	return Default()
}

// New initializes a new logger with level and provided outputs.
func New(level Level, outputs ...Output) Logger {
	l := logrus.New()
	l.SetLevel(level.logrus())
	l.SetNoLock()
	l.SetReportCaller(true)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006/01/02 15:04:05",
		DataKey:         "fields",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
		},
	})

	if len(outputs) == 0 {
		outputs = append(outputs, OutputStd())
	}

	l.SetOutput(&outputCluster{outputs})

	log := &logger{
		entry: l.WithContext(context.Background()),
	}

	return log
}

type logger struct {
	entry *logrus.Entry
}

func (l *logger) GetLevel() Level {
	return newLevelFromLogrus(l.entry.Logger.Level)
}

func (l *logger) SetOutput(output io.Writer) {
	l.entry.Logger.SetOutput(output)
}

func (l *logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *logger) Copy() Logger {
	copied := l.entry.Dup()
	return &logger{
		entry: copied,
	}
}

func (l *logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logger) Log(level Level, args ...interface{}) {
	l.entry.Log(level.logrus(), args...)
}

func (l *logger) Logf(level Level, format string, args ...interface{}) {
	l.entry.Logf(level.logrus(), format, args...)
}

func (l *logger) Logln(level Level, args ...interface{}) {
	l.entry.Logln(level.logrus(), args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.entry.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logger) Panicln(args ...interface{}) {
	l.entry.Panicln(args...)
}

func (l *logger) Print(args ...interface{}) {
	l.entry.Print(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.entry.Printf(format, args...)
}

func (l *logger) Println(args ...interface{}) {
	l.entry.Println(args...)
}

func (l *logger) String() (string, error) {
	return l.entry.String()
}

func (l *logger) Trace(args ...interface{}) {
	l.entry.Trace(args...)
}

func (l *logger) Tracef(format string, args ...interface{}) {
	l.entry.Tracef(format, args...)
}

func (l *logger) Traceln(args ...interface{}) {
	l.entry.Traceln(args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.entry.Warning(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.entry.Warningf(format, args...)
}

func (l *logger) Warningln(args ...interface{}) {
	l.entry.Warningln(args...)
}

func (l *logger) Warnln(args ...interface{}) {
	l.entry.Warnln(args...)
}

func (l *logger) WithContext(ctx context.Context) Logger {
	return &logger{
		entry: l.entry.WithContext(ctx),
	}
}

func (l *logger) WithError(err error) Logger {
	if err == nil {
		return l
	}

	return &logger{
		entry: l.entry.WithField("error", fmt.Sprintf("%+v", err)),
	}
}

func (l *logger) WithField(k string, v interface{}) Logger {
	return &logger{
		entry: l.entry.WithField(k, v),
	}
}

func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

func (l *logger) WithTime(t time.Time) Logger {
	return &logger{
		entry: l.entry.WithTime(t),
	}
}

func (l *logger) Writer() *io.PipeWriter {
	return l.entry.Writer()
}

func (l *logger) WriterLevel(level Level) *io.PipeWriter {
	return l.entry.WriterLevel(level.logrus())
}

func (l *logger) WithFunc(function string) Logger {
	return &logger{
		entry: l.entry.WithField("func", function),
	}
}

func (l *logger) Attach(ctx context.Context) (context.Context, Logger) {
	ll := l.Copy()
	return context.WithValue(ctx, logKey{}, ll), ll
}
