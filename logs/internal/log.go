package internal

import (
	"sync/atomic"
)

var (
	// DefaultLogger is the default logger.
	DefaultLogger atomic.Value

	DefaultTimeFormat atomic.Value
)

const (
	DefaultFormat = "2006/01/02 15:04:05"
)
