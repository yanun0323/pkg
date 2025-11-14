package ws

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yanun0323/errors"
	"github.com/yanun0323/pkg/channel"
	"github.com/yanun0323/pkg/sys"
)

var (
	DefaultWaitingMessageTimeout = 15 * time.Second
)

type Producer[T any] interface {
	Produce() <-chan T
}

type SubscriberID int64
type Subscriber[T any] func(T)

type Publisher[T any] struct {
	producer Producer[T]

	subsMu     sync.RWMutex
	subs       map[SubscriberID]chan T
	subsNextID SubscriberID

	stop context.CancelFunc

	start atomic.Bool
	end   atomic.Bool
}

const (
	_defaultSubscriberCount = 1000
	_defaultSubscriberCap   = 1000
)

func NewPublisher[T any](producer Producer[T]) *Publisher[T] {
	return &Publisher[T]{
		producer: producer,
		subs:     make(map[SubscriberID]chan T, _defaultSubscriberCount),
	}
}

func (pub *Publisher[T]) Len() int {
	pub.subsMu.RLock()
	defer pub.subsMu.RUnlock()

	return len(pub.subs)
}

func (pub *Publisher[T]) Start(ctx context.Context) {
	if pub.start.Swap(true) {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	pub.stop = cancel

	go pub.consumeMessage(ctx)
}

func (pub *Publisher[T]) consumeMessage(ctx context.Context) {
	for {
		select {
		case <-sys.Shutdown():
			return
		case <-ctx.Done():
			return
		case msg, ok := <-pub.producer.Produce():
			if !ok {
				pub.subsMu.Lock()
				pub.end.Store(true)
				for _, sub := range pub.subs {
					channel.SafeClose(sub)
				}
				pub.subsMu.Unlock()

				return
			}

			pub.subsMu.RLock()
			for id, sub := range pub.subs {
				if ok := channel.TryPush(sub, msg); !ok {
					fmt.Printf("message dropped! %d subscriber channel is full\n", id)
				}
			}
			pub.subsMu.RUnlock()
		}
	}
}

func (pub *Publisher[T]) Subscribe(ctx context.Context, sub Subscriber[T]) (unsubscribe func()) {
	ch := make(chan T, _defaultSubscriberCap)
	pub.subsMu.Lock()
	defer pub.subsMu.Unlock()

	if pub.end.Load() {
		return
	}

	id := pub.subsNextID
	pub.subsNextID++
	pub.subs[id] = ch

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer channel.SafeClose(ch)

		for {
			select {
			case <-sys.Shutdown():
				return
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}

				sub(msg)
			}
		}
	}()

	return func() {
		pub.subsMu.Lock()
		defer pub.subsMu.Unlock()
		delete(pub.subs, id)

		cancel()
	}
}

func SubscribeAndWait[T any, Result any](ctx context.Context, pub *Publisher[T], mapping func(T) (Result, bool), before func(context.Context) error, isWaitTarget func(context.Context, Result) bool, timeout ...time.Duration) (Result, error) {
	done := make(chan error, 1)
	defer channel.SafeClose(done)

	msg := make(chan struct{})
	defer channel.SafeClose(msg)

	waitTimeout := DefaultWaitingMessageTimeout
	if len(timeout) != 0 && timeout[0] > 0 {
		waitTimeout = timeout[0]
	}

	ctx, cancel := context.WithTimeout(ctx, waitTimeout)
	defer cancel()

	var result Result
	unsubscribe := pub.Subscribe(ctx, func(t T) {
		r, ok := mapping(t)
		if !ok {
			return
		}

		if isWaitTarget(ctx, r) {
			result = r
			channel.SafeClose(msg)
		}
	})
	defer unsubscribe()

	if err := before(ctx); err != nil {
		return result, errors.Wrap(err, "execute before function")
	}

	select {
	case <-sys.Shutdown():
		done <- context.Canceled
	case <-ctx.Done():
		done <- ctx.Err()
	case <-msg:
		done <- nil
	}

	return result, errors.Wrap(<-done, "websocket message")
}
