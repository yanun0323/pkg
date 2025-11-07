package sys

import (
	"os"
	"os/signal"
	"syscall"
)

func Shutdown() <-chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}
