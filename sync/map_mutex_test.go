package sync

import "testing"

func TestMapMutex(t *testing.T) {
	sm := &MapMutex[string]{}

	{
		sm.Lock("test")
		ok := sm.TryLock("test")
		if ok {
			t.Errorf("TryLock should return false")
		}
		sm.Unlock("test")
	}

	{
		if ok := sm.TryLock("test"); !ok {
			t.Errorf("TryLock should return true")
		}
		sm.Unlock("test")
	}

	{
		sm.Lock("test")
		sm.Lock("test2")
		sm.Unlock("test")
		sm.Unlock("test2")
	}
}
