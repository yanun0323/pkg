package syncs

import "sync"

type MapMutex[Key comparable] struct {
	_ noCopy

	mu   sync.Mutex
	data map[Key]*sync.Mutex
}

func (m *MapMutex[Key]) init(key Key) {
	if m.data == nil {
		m.data = make(map[Key]*sync.Mutex)
	}

	if m.data[key] == nil {
		m.data[key] = &sync.Mutex{}
	}
}

func (m *MapMutex[Key]) Lock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].Lock()
}

func (m *MapMutex[Key]) TryLock(key Key) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	return m.data[key].TryLock()
}

func (m *MapMutex[Key]) Unlock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].Unlock()
}
