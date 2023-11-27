package logs

import (
	"fmt"
	"sync"
	"testing"
)

func TestWithField(t *testing.T) {
	wg := sync.WaitGroup{}
	l := New("test", LevelDebug)
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			funcName := fmt.Sprintf("fund-%d", i)
			ll := l.WithField("func", funcName)
			ll.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}

func TestWithFields(t *testing.T) {
	wg := sync.WaitGroup{}
	l := New(LevelDebug)
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			funcName := fmt.Sprintf("fund-%d", i)
			l = l.WithFields(map[string]interface{}{"func": funcName})
			l.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}
