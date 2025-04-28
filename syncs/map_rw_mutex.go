package syncs

import "sync"

type MapRWMutex[Key comparable] struct {
	_ noCopy

	mu   sync.Mutex
	data map[Key]*sync.RWMutex
}

func (m *MapRWMutex[Key]) init(key Key) {
	if m.data == nil {
		m.data = make(map[Key]*sync.RWMutex)
	}

	if m.data[key] == nil {
		m.data[key] = &sync.RWMutex{}
	}
}

func (m *MapRWMutex[Key]) Lock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].Lock()
}

func (m *MapRWMutex[Key]) TryLock(key Key) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	return m.data[key].TryLock()
}

func (m *MapRWMutex[Key]) Unlock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].Unlock()
}

func (m *MapRWMutex[Key]) RLock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].RLock()
}

func (m *MapRWMutex[Key]) TryRLock(key Key) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	return m.data[key].TryRLock()
}

func (m *MapRWMutex[Key]) RUnlock(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.init(key)
	m.data[key].RUnlock()
}
