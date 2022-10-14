package logs

import (
	"bytes"
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

/*
Get stack information form error.
*/
func GetStack(err error) string {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	sterr, ok := err.(stackTracer)
	if !ok {
		err = errors.WithStack(err)
		sterr = err.(stackTracer)
	}
	st := sterr.StackTrace()

	s := bytes.Buffer{}
	for i := 0; i < len(st); i++ {
		f := st[i]
		fileName := getFile(f)
		line := getLine(f)
		s.Write([]byte(fmt.Sprintf("%s:%d", fileName, line)))
		if i < len(st)-1 {
			s.WriteByte('\n')
		}
	}

	return s.String()
}

func getFile(f errors.Frame) string {
	pc := uintptr(f) - 1
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(pc)
	return file
}

func getLine(f errors.Frame) int {
	pc := uintptr(f) - 1
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(pc)
	return line
}
