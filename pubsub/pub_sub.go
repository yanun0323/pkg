package pubsub

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
	DefaultSubscriberCap         = 1000
	DefaultSubscriberMessageCap  = 1000
)

type Producer[T any] interface {
	Start(context.Context)
	Produce() <-chan T
}

type SubscriberID int64
type Subscriber[T any] func(T)

type Publisher[P Producer[T], T any] struct {
	producer P

	subsMu     sync.RWMutex
	subs       map[SubscriberID]chan T
	subsNextID SubscriberID

	stop context.CancelFunc

	start atomic.Bool
	end   atomic.Bool
}

func NewPublisher[P Producer[T], T any](producer P, subscribeCap ...int) *Publisher[P, T] {
	caps := DefaultSubscriberCap
	if len(subscribeCap) != 0 && subscribeCap[0] > 0 {
		caps = subscribeCap[0]
	}

	return &Publisher[P, T]{
		producer: producer,
		subs:     make(map[SubscriberID]chan T, caps),
	}
}

func (pub *Publisher[P, T]) Producer() P {
	return pub.producer
}

func (pub *Publisher[P, T]) Len() int {
	pub.subsMu.RLock()
	defer pub.subsMu.RUnlock()

	return len(pub.subs)
}

func (pub *Publisher[P, T]) Start(ctx context.Context) {
	if pub.start.Swap(true) {
		return
	}

	ctx, cancel := context.WithCancel(ctx)
	pub.stop = cancel

	go pub.consumeMessage(ctx)
	pub.producer.Start(ctx)
}

func (pub *Publisher[P, T]) consumeMessage(ctx context.Context) {
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

func (pub *Publisher[P, T]) Subscribe(ctx context.Context, sub Subscriber[T], messageCap ...int) (unsubscribe func()) {
	caps := DefaultSubscriberMessageCap
	if len(messageCap) != 0 && messageCap[0] > 0 {
		caps = messageCap[0]
	}

	ch := make(chan T, caps)
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

func (pub *Publisher[P, T]) SubscribeAndWait(ctx context.Context, send func(context.Context, P) error, isExpected func(context.Context, T) bool, timeout ...time.Duration) error {
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

	unsubscribe := pub.Subscribe(ctx, func(t T) {
		if isExpected(ctx, t) {
			channel.SafeClose(msg)
		}
	})
	defer unsubscribe()

	if err := send(ctx, pub.producer); err != nil {
		return errors.Wrap(err, "execute before function")
	}

	select {
	case <-sys.Shutdown():
		done <- context.Canceled
	case <-ctx.Done():
		done <- ctx.Err()
	case <-msg:
		done <- nil
	}

	if err := <-done; err != nil {
		return errors.Wrap(err, "unexpected result")
	}

	return nil
}
