package logs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/buger/jsonparser"
)

type Output io.Writer

// OutputStd return a std output.
func OutputStd() Output {
	return &stdout{}
}

// OutputFile return an file output.
func OutputFile(dir, filename string) Output {
	w, _ := os.OpenFile(fmt.Sprintf("%s/%s.log", getAbsPath(dir), filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return w
}

type outputContainer struct {
	writers []io.Writer
}

var (
	errorReplaceString = []string{"\\n", "\n", "\\t", "\t"}
)

const (
	errorMsgBegin = "\n"
)

func newOutputContainer(ops ...Output) io.Writer {
	writers := make([]io.Writer, 0, len(ops))
	for _, op := range ops {
		writers = append(writers, op)
	}
	return &outputContainer{
		writers: writers,
	}
}

func (o *outputContainer) Write(p []byte) (int, error) {
	for _, w := range o.writers {
		if _, err := w.Write(p); err != nil {
			fmt.Println(err)
		}
	}

	return len(p), nil
}

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite
)

const (
	colorBrightBlack = iota + 90
	colorBrightRed
	colorBrightGreen
	colorBrightYellow
	colorBrightBlue
	colorBrightMagenta
	colorBrightCyan
	colorBrightWhite
)

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

type stdout struct{}

func (*stdout) Write(p []byte) (int, error) {
	buf := bytes.Buffer{}

	timestamp, _ := jsonparser.GetString(p, "@timestamp")
	level, _ := jsonparser.GetString(p, "level")

	msg, _ := jsonparser.GetString(p, "msg")

	level = strings.ToUpper(level)
	buf.WriteString(colorize(timestamp, colorBlack))
	buf.WriteByte(' ')
	buf.WriteString(colorize(getTitle(level), getLevelColor(level)))
	buf.WriteByte(' ')
	buf.WriteString(msg)

	fields := [][]byte{}
	errorMsg, _ := jsonparser.GetString(p, "fields", "error")

	if errorMsg != "" {
		errorMsg = errorMsgBegin + colorize("[error stack] ", colorBrightRed) + errorMsg
		errorMsgReplacer := strings.NewReplacer(errorReplaceString...)
		errorMsg = errorMsgReplacer.Replace(errorMsg)
		p = jsonparser.Delete(p, "fields", "error")
	}

	jsonparser.ObjectEach(p, func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		field := bytes.Buffer{}
		field.WriteString(colorize("["+string(key)+"] ", colorMagenta))
		field.WriteString(colorize(string(value), colorBlack))
		fields = append(fields, field.Bytes())

		return nil
	}, "fields")

	if len(fields) != 0 {
		buf.WriteString("  ")
		buf.Write(bytes.Join(fields, []byte(" ")))
	}

	if len(errorMsg) != 0 {
		buf.WriteString(errorMsg)
	}

	buf.WriteString("\n")

	fmt.Fprint(os.Stdout, buf.String())

	return len(p), nil
}

func getTitle(level string) string {
	if level == "WARNING" {
		level = "WARN"
	}
	// span := 5 - len(level)
	return level
}

func getLevelColor(level string) int {
	switch level {
	case "DEBUG":
		return colorBlue
	case "INFO":
		return colorGreen
	case "ERROR", "PANIC":
		return colorRed + 10
	case "WARNING":
		return colorYellow
	case "FATAL":
		return colorBrightRed
	default:
		return colorBlue
	}
}
