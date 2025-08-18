package atomic

import "sync/atomic"

type Value[T any] struct {
	_ noCopy

	v atomic.Value
}

func (v *Value[T]) Load() T {
	val, _ := v.v.Load().(T)
	return val
}

func (v *Value[T]) Store(val T) {
	v.v.Store(val)
}

func (v *Value[T]) Swap(val T) (T, bool) {
	old, ok := v.v.Swap(val).(T)
	return old, ok
}

func (v *Value[T]) CompareAndSwap(old, new T) bool {
	return v.v.CompareAndSwap(old, new)
}
