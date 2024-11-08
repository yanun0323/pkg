package pubsub_test

import (
	"context"
	"testing"

	"github.com/yanun0323/pkg/pubsub"
)

func TestProducer(t *testing.T) {
	p := pubsub.NewProducer[int]()

	receiver := make(chan int, 10)
	p.Subscribe(func(ctx context.Context, val int, close bool) (keep bool) {
		receiver <- val
		return true
	})

}
