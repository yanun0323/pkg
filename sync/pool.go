package sync

import "sync"

type resettable interface {
	Reset()
}

type Pool[T resettable] interface {
	Get() T
	Put(T)
}

type pool[T resettable] struct {
	sync.Pool
}

func NewPool[T resettable](newFn func() T) Pool[T] {
	return &pool[T]{
		Pool: sync.Pool{New: func() any { return newFn() }},
	}
}

func (p *pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *pool[T]) Put(t T) {
	t.Reset()
	p.Pool.Put(t)
}
