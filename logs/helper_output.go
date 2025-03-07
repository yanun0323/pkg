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

type outputCluster struct {
	writers []Output
}

var (
	errorReplaceString = []string{"\\n", "\n", "\\t", "\t"}
	errorMsgReplacer   = strings.NewReplacer(errorReplaceString...)
)

const (
	errorMsgBegin = "\n"
)

func (o *outputCluster) Write(p []byte) (int, error) {
	for _, w := range o.writers {
		if _, err := w.Write(p); err != nil {
			fmt.Println(err)
		}
	}

	return len(p), nil
}

type stdout struct{}

func (*stdout) Write(p []byte) (int, error) {
	buf := bytes.Buffer{}

	timestamp, _ := jsonparser.GetString(p, "timestamp")
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
