package pubsub

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Subscribable[T any] interface {
	Subscribe(...Consumer[T])
}

type Producer[T any] interface {
	// Publish publishes a value to all subscribers.
	//
	// It uses the default timeout (15 seconds) if the constructor not provided.
	Publish(context.Context, T)

	// Produce runs a goroutine to publish the value to the subscribers.
	Produce(context.Context, ...<-chan T)

	// Subscribe subscribes a consumer to the producer.
	Subscribe(...Consumer[T])

	// Clear unsubscribes all consumer from the producer.
	Clear()
}

const _defaultTimeout = time.Second * 15

// Consumer is a function that consumes a value from the producer.
//
// The function should return true to keep consuming the value, or false to stop consuming.
type Consumer[T any] func(ctx context.Context, val T, close bool) (keep bool)

type producer[T any] struct {
	mu          sync.RWMutex
	subscribers map[string]Consumer[T]
	timeout     time.Duration
}

// NewProducer creates a new producer.
//
// use default message publishing timeout (15 seconds) if not provided.
func NewProducer[T any](messagePublishTimeout ...time.Duration) Producer[T] {
	timeout := _defaultTimeout
	if len(messagePublishTimeout) != 0 {
		timeout = messagePublishTimeout[0]
	}
	return &producer[T]{
		mu:          sync.RWMutex{},
		subscribers: map[string]Consumer[T]{},
		timeout:     timeout,
	}
}

func (p *producer[T]) Publish(ctx context.Context, t T) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(len(p.subscribers))
	for k, f := range p.subscribers {
		if f == nil {
			go p.Unsubscribe(k)
			continue
		}

		go func(k string, f Consumer[T]) {
			wg.Done()
			keep := f(ctx, t, false)
			if !keep {
				go p.Unsubscribe(k)
			}
		}(k, f)
	}
	wg.Wait()
}

func (p *producer[T]) Unsubscribe(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribers, key)
}

func (p *producer[T]) Produce(ctx context.Context, pipes ...<-chan T) {
	for _, pipe := range pipes {
		go func(pp <-chan T) {
			for {
				select {
				case <-ctx.Done():
					p.Clear()
					return
				case t, ok := <-pp:
					if !ok {
						return
					}

					p.Publish(ctx, t)
				}
			}
		}(pipe)
	}
}

func (p *producer[T]) Subscribe(consumers ...Consumer[T]) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, f := range consumers {
		if f == nil {
			continue
		}

		p.subscribers[uuid.New().String()] = f
	}
}

func (p *producer[T]) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, f := range p.subscribers {
		if f == nil {
			continue
		}

		go func(f Consumer[T]) {
			ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
			defer cancel()

			_ = f(ctx, *new(T), true)
		}(f)
	}

	p.subscribers = map[string]Consumer[T]{}
}
