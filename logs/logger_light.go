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
	ctx        *value[context.Context]
	level      *value[Level]
	fields     *value[map[string]any]
	output     *value[Output]
	timeFormat *value[string]
	time       *value[*time.Time]
}

func NewLightweightLogger(level Level, outputs ...Output) Logger {
	return &loggerLight{
		ctx:        newValue(context.Background()),
		level:      newValue(LevelInfo),
		fields:     newValue(map[string]any{}),
		output:     newValue[Output](&outputCluster{outputs}),
		timeFormat: newValue(_defaultTimestampFormat),
		time:       newValue[*time.Time](),
	}
}

func (l *loggerLight) copy() *loggerLight {
	var copied map[string]any
	l.fields.Update(func(fields map[string]any) map[string]any {
		copied = make(map[string]any, len(fields))
		for k, v := range fields {
			copied[k] = v
		}
		return fields
	})

	return &loggerLight{
		ctx:        l.ctx.Copy(),
		level:      l.level.Copy(),
		fields:     newValue(copied),
		output:     l.output.Copy(),
		timeFormat: l.timeFormat.Copy(),
		time:       l.time.Copy(),
	}
}

func (l *loggerLight) Copy() Logger {
	return l.copy()
}

func (l *loggerLight) GetLevel() Level {
	return l.level.Load()
}

func (l *loggerLight) SetOutput(output io.Writer) {
	l.output.Store(output)
}

func (l *loggerLight) SetTimestampFormat(format string) {
	l.timeFormat.Store(format)
}

func (l *loggerLight) WithContext(ctx context.Context) Logger {
	ll := l.copy()
	ll.ctx.Store(ctx)
	return ll
}

func (l *loggerLight) WithTime(t time.Time) Logger {
	ll := l.copy()
	ll.time.Store(&t)
	return ll
}

func (l *loggerLight) WithField(key string, value any) Logger {
	ll := l.copy()
	ll.fields.Update(func(fields map[string]any) map[string]any {
		fields[key] = value
		return fields
	})
	return ll
}

func (l *loggerLight) WithFields(fields map[string]any) Logger {
	ll := l.copy()
	ll.fields.Update(func(fields map[string]any) map[string]any {
		for k, v := range fields {
			fields[k] = v
		}
		return fields
	})
	return ll
}

func (l *loggerLight) WithFunc(function string) Logger {
	ll := l.copy()
	ll.fields.Update(func(fields map[string]any) map[string]any {
		fields["func"] = function
		return fields
	})
	return ll
}

func (l *loggerLight) WithError(err error) Logger {
	ll := l.copy()
	ll.fields.Update(func(fields map[string]any) map[string]any {
		fields["error"] = err
		return fields
	})
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
	if tt := o.l.time.Load(); tt != nil {
		t = *tt
	}
	buf := bytes.NewBuffer(nil)

	buf.WriteString(colorize(t.Format(o.l.timeFormat.Load()), colorBlack))
	buf.WriteByte(' ')
	buf.WriteString(colorize(getTitle(o.level.String()), getLevelColor(o.level.String())))
	buf.WriteByte(' ')
	buf.Write(msg)

	fields := o.l.fields.Load()
	err, errOK := fields["error"]
	delete(fields, "error")

	for k, v := range fields {
		buf.WriteString("  ")
		buf.WriteString(colorize("["+k+"] ", colorMagenta))
		buf.WriteString(colorize(fmt.Sprintf("%v", v), colorBlack))
	}

	if ctx := o.l.ctx.Load(); ctx != nil {
		buf.WriteString("  ")
		buf.WriteString(colorize("[context] ", colorMagenta))
		buf.WriteString(colorize(fmt.Sprintf("%v", ctx), colorBlack))
	}

	if errOK {
		buf.WriteByte('\n')
		buf.WriteString(colorize("[error stack] ", colorBrightRed))
		buf.WriteString(errorMsgReplacer.Replace(fmt.Sprintf("%+v", err)))
	}

	buf.WriteString("\n")

	p := buf.String()
	fmt.Fprint(o.l.output.Load(), p)

	return len(p), nil
}
