package syncs

import "sync"

type Map[K comparable, V any] struct {
	_ noCopy

	sync.Map
}

func (m *Map[K, V]) Clear() {
	m.Map.Clear()
}
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.Map.CompareAndDelete(key, old)
}

func (m *Map[K, V]) CompareAndSwap(key K, old V, new V) (swapped bool) {
	return m.Map.CompareAndSwap(key, old, new)
}

func (m *Map[K, V]) Delete(key K) {
	m.Map.Delete(key)
}

func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.Map.Load(key)
	if !ok {
		return value, false
	}

	value, ok = v.(V)
	return value, ok
}

func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.Map.LoadAndDelete(key)
	if !loaded {
		return value, false
	}

	value, ok := v.(V)
	return value, ok
}

func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.Map.LoadOrStore(key, value)
	if !loaded {
		return value, false
	}

	actual, ok := v.(V)
	return actual, ok
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.Map.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) Store(key K, value V) {
	m.Map.Store(key, value)
}
