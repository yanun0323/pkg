package sys

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	shutdown = make(chan os.Signal)
	once     sync.Once
)

func Shutdown() <-chan os.Signal {
	once.Do(func() {
		go func() {
			agg := make(chan os.Signal, 1)
			signal.Notify(agg, syscall.SIGINT, syscall.SIGTERM)
			<-agg
			println()
			close(shutdown)
		}()
	})

	return shutdown
}

