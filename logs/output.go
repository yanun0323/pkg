package logs

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/buger/jsonparser"
)

type output struct {
	app     string
	writers []io.Writer
}

var (
	errorReplaceString = []string{"\\n", "\n", "\\t", "\t"}
)

const (
	errorMsgBegin = "\n"
)

func (o *output) new(outs, app string) io.Writer {
	writers := []io.Writer{}

	if strings.Contains(outs, "stdout") {
		writers = append(writers, &stdout{})
	}

	if strings.Contains(outs, "file") {
		w, _ := os.OpenFile(fmt.Sprintf("%s/%s.log", GetAbsPath(app, "log"), fmt.Sprintf("/%s", app)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		writers = append(writers, &file{
			writer: w,
		})
	}

	return &output{
		writers: writers,
	}
}

func (o *output) Write(p []byte) (int, error) {
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

const (
	styleBold = 1
)

func colorize(s interface{}, c int) string {
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

type stdout struct{}

func (*stdout) Write(p []byte) (int, error) {
	buf := bytes.Buffer{}

	datatime, _ := jsonparser.GetString(p, "datatime")
	level, _ := jsonparser.GetString(p, "level")
	// eventID, _, _, _ := jsonparser.Get(p, "fields", "eventId")
	msg, _ := jsonparser.GetString(p, "msg")
	file, _ := jsonparser.GetString(p, "file")
	line := file[strings.LastIndex(file, "/")+1:]

	level = strings.ToUpper(level)
	buf.WriteString(colorize(datatime, colorBrightBlack))
	buf.WriteString(" ")
	buf.WriteString(colorize(getTitle(level), getLevelColor(level)))
	buf.WriteString(" ")
	// buf.WriteString("[")
	// buf.Write([]byte(colorize(string(eventID), colorGreen)))
	// buf.WriteString("]")
	buf.WriteString(" ")
	buf.WriteString(msg)
	buf.WriteString(" ")
	buf.WriteString("@")
	buf.WriteString(line)

	fields := [][]byte{}
	errorMsg, _ := jsonparser.GetString(p, "fields", "error")

	if errorMsg != "" {
		errorMsg = errorMsgBegin + errorMsg
		errorMsgReplacer := strings.NewReplacer(errorReplaceString...)
		errorMsg = errorMsgReplacer.Replace(errorMsg)
		p = jsonparser.Delete(p, "fields", "error")
	}

	jsonparser.ObjectEach(p, func(key []byte, value []byte, _ jsonparser.ValueType, _ int) error {
		field := bytes.Buffer{}
		field.Write([]byte(colorize(string(key), colorCyan)))
		field.Write([]byte("="))
		field.Write(value)
		fields = append(fields, field.Bytes())

		return nil
	}, "fields")

	if len(fields) > 0 {
		buf.WriteString(" ")
		buf.WriteString("[")
		buf.Write(bytes.Join(fields, []byte(",")))
		buf.WriteString("]")
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
	span := 6 - len(level)
	return " " + level + strings.Repeat(" ", span)
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

type file struct {
	writer *os.File
}

func (f *file) Write(p []byte) (int, error) {
	fmt.Fprint(f.writer, string(p))

	return len(p), nil
}
