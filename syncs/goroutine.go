package syncs

import (
	"context"
	"log"
	"time"
)

func Go(fn func(), timeout ...time.Duration) {
	if len(timeout) == 0 {
		go func() {
			defer func() {
				if err, ok := recover().(error); ok && err != nil {
					log.Printf("\x1b[41mERROR\x1b[0m %+v", err)
				}
			}()

			fn()
		}()

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout[0])
	defer cancel()

	signal := make(chan struct{}, 1)
	defer close(signal)

	go func() {
		defer func() {
			if err, ok := recover().(error); ok && err != nil {
				log.Printf("\x1b[41mERROR\x1b[0m %+v", err)
			}
		}()

		fn()

		select {
		case <-ctx.Done():
			return
		default:
			signal <- struct{}{}
		}
	}()

	select {
	case <-signal:
	case <-ctx.Done():
	}
}
