package logs

import "sync"

type value[T any] struct {
	mu  sync.RWMutex
	val T
}

func newValue[T any](val ...T) *value[T] {
	var v T
	if len(val) != 0 {
		v = val[0]
	}

	return &value[T]{
		mu:  sync.RWMutex{},
		val: v,
	}
}

func (v *value[T]) Load() T {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.val
}

func (v *value[T]) Store(val T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.val = val
}

func (v *value[T]) Swap(val T) T {
	v.mu.Lock()
	defer v.mu.Unlock()

	old := v.val
	v.val = val

	return old
}

func (v *value[T]) Update(fn func(T) T) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.val = fn(v.val)
}

func (v *value[T]) Copy() *value[T] {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return newValue(v.val)
}
