package pubsub_test

import (
	"context"
	"testing"
	"time"

	"github.com/yanun0323/pkg/pubsub"
)

func TestProducer(t *testing.T) {
	p := pubsub.NewProducer[int]()

	receiver := make(chan int, 10)
	p.Subscribe(func(ctx context.Context, val int, close bool) (keep bool) {
		receiver <- val
		return true
	})

	p.Publish(context.Background(), 1)

	select {
	case <-time.After(3 * time.Second):
		t.Fatal("timeout")
	case <-receiver:
	}
}
