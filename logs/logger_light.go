package logs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"
)

const (
	_defaultTimestampFormat = "2006/01/02 15:04:05"
)

type loggerLight struct {
	ctx        context.Context
	level      Level
	fields     map[string]any
	output     Output
	timeFormat string
	time       *time.Time
}

func NewLightweightLogger(level Level, outputs ...Output) Logger {
	return &loggerLight{
		ctx:        context.Background(),
		level:      LevelInfo,
		fields:     map[string]any{},
		output:     &outputCluster{outputs},
		timeFormat: _defaultTimestampFormat,
		time:       nil,
	}
}

func (l *loggerLight) copy() *loggerLight {
	copied := make(map[string]any, len(l.fields))
	for k, v := range l.fields {
		copied[k] = v
	}

	return &loggerLight{
		ctx:        l.ctx,
		level:      l.level,
		fields:     copied,
		output:     l.output,
		timeFormat: l.timeFormat,
		time:       l.time,
	}
}

func (l *loggerLight) Copy() Logger {
	return l.copy()
}

func (l *loggerLight) GetLevel() Level {
	return l.level
}

func (l *loggerLight) SetOutput(output io.Writer) {
	l.output = output
}

func (l *loggerLight) SetTimestampFormat(format string) {
	l.timeFormat = format
}

func (l *loggerLight) WithContext(ctx context.Context) Logger {
	ll := l.copy()
	ll.ctx = ctx
	return ll
}

func (l *loggerLight) WithTime(t time.Time) Logger {
	ll := l.copy()
	ll.time = &t
	return ll
}

func (l *loggerLight) WithField(key string, value any) Logger {
	ll := l.copy()
	ll.fields[key] = value
	return ll
}

func (l *loggerLight) WithFields(fields map[string]any) Logger {
	ll := l.copy()
	for k, v := range fields {
		ll.fields[k] = v
	}
	return ll
}

func (l *loggerLight) WithFunc(function string) Logger {
	ll := l.copy()
	ll.fields["func"] = function
	return ll
}

func (l *loggerLight) WithError(err error) Logger {
	ll := l.copy()
	ll.fields["error"] = err
	return ll
}

func (l *loggerLight) Writer() *io.PipeWriter {
	// BUG: Fix me
	return nil
}

func (l *loggerLight) WriterLevel(level Level) *io.PipeWriter {
	// BUG: Fix me
	return nil
}

func (l *loggerLight) Attach(ctx context.Context) (context.Context, Logger) {
	ll := l.Copy()
	return context.WithValue(ctx, logKey{}, ll), ll
}

func (l *loggerLight) Log(level Level, args ...interface{}) {
	l.print(level, args...)
}

func (l *loggerLight) Logf(level Level, format string, args ...interface{}) {
	l.printf(level, format, args...)
}

func (l *loggerLight) Trace(args ...interface{}) {
	l.print(LevelTrace, args...)
}

func (l *loggerLight) Tracef(format string, args ...interface{}) {
	l.printf(LevelTrace, format, args...)
}

func (l *loggerLight) Debug(args ...interface{}) {
	l.print(LevelDebug, args...)
}

func (l *loggerLight) Debugf(format string, args ...interface{}) {
	l.printf(LevelDebug, format, args...)
}

func (l *loggerLight) Info(args ...interface{}) {
	l.print(LevelInfo, args...)
}

func (l *loggerLight) Infof(format string, args ...interface{}) {
	l.printf(LevelInfo, format, args...)
}

func (l *loggerLight) Warn(args ...interface{}) {
	l.print(LevelWarn, args...)
}

func (l *loggerLight) Warnf(format string, args ...interface{}) {
	l.printf(LevelWarn, format, args...)
}

func (l *loggerLight) Warning(args ...interface{}) {
	l.print(LevelWarn, args...)
}

func (l *loggerLight) Warningf(format string, args ...interface{}) {
	l.printf(LevelWarn, format, args...)
}

func (l *loggerLight) Error(args ...interface{}) {
	l.print(LevelError, args...)
}

func (l *loggerLight) Errorf(format string, args ...interface{}) {
	l.printf(LevelError, format, args...)
}

func (l *loggerLight) Fatal(args ...interface{}) {
	l.print(LevelFatal, args...)
}

func (l *loggerLight) Fatalf(format string, args ...interface{}) {
	l.printf(LevelFatal, format, args...)
}

func (l *loggerLight) Panic(args ...interface{}) {
	l.print(LevelPanic, args...)
}

func (l *loggerLight) Panicf(format string, args ...interface{}) {
	l.printf(LevelPanic, format, args...)
}

func (l *loggerLight) print(level Level, args ...interface{}) {
	if !l.GetLevel().isAvailable(level) {
		return
	}

	lightOutput{l, level}.Write([]byte(fmt.Sprint(args...)))
}

func (l *loggerLight) printf(level Level, format string, args ...interface{}) {
	if !l.GetLevel().isAvailable(level) {
		return
	}

	lightOutput{l, level}.Write([]byte(fmt.Sprintf(format, args...)))
}

type lightOutput struct {
	l     *loggerLight
	level Level
}

func (o lightOutput) Write(msg []byte) (int, error) {
	t := time.Now()
	if o.l.time != nil {
		t = *o.l.time
	}
	buf := bytes.NewBuffer(nil)

	buf.WriteString(colorize(t.Format(o.l.timeFormat), colorBlack))
	buf.WriteByte(' ')
	buf.WriteString(colorize(getTitle(o.level.String()), getLevelColor(o.level.String())))
	buf.WriteByte(' ')
	buf.Write(msg)

	for k, v := range o.l.fields {
		if k == "error" {
			continue
		}
		buf.WriteString("  ")
		buf.WriteString(colorize("["+k+"] ", colorMagenta))
		buf.WriteString(colorize(fmt.Sprintf("%v", v), colorBlack))
	}

	if o.l.ctx != nil {
		buf.WriteString("  ")
		buf.WriteString(colorize("[context] ", colorMagenta))
		buf.WriteString(colorize(fmt.Sprintf("%v", o.l.ctx), colorBlack))
	}

	if err, errOK := o.l.fields["error"]; errOK && err != nil {
		buf.WriteByte('\n')
		buf.WriteString(colorize("[error stack] ", colorBrightRed))
		buf.WriteString(errorMsgReplacer.Replace(fmt.Sprintf("%+v", err)))
	}

	buf.WriteString("\n")

	p := buf.String()
	fmt.Fprint(o.l.output, p)

	return len(p), nil
}
