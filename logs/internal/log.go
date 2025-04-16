package internal

import (
	"sync/atomic"
)

const (
	_defaultFormat = "2006/01/02 15:04:05"
)

var (
	defaultTimeFormat atomic.Value
)

func GetDefaultTimeFormat() string {
	s, ok := defaultTimeFormat.Load().(string)
	if !ok || len(s) == 0 {
		return _defaultFormat
	}

	return s
}

func SetDefaultTimeFormat(format string) {
	if len(format) != 0 {
		defaultTimeFormat.Store(format)
	}
}

func NewValue(val any) (v atomic.Value) {
	v.Store(val)
	return
}
