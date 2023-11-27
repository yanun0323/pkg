package logs

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestWithField(t *testing.T) {
	wg := sync.WaitGroup{}
	l := newWithOutput("test", LevelDebug, "stdout")
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			// ctx := context.Background()
			funcName := fmt.Sprintf("fund-%d", i)
			ll := l.WithField("func", funcName)
			ll.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}

func TestWithFields(t *testing.T) {
	wg := sync.WaitGroup{}
	l := newWithOutput("test", LevelDebug, "stdout")
	count := 10
	wg.Add(count)
	for i := 1; i <= count; i++ {
		go func(i int) {
			defer wg.Done()
			ctx := context.Background()
			funcName := fmt.Sprintf("fund-%d", i)
			l, _ = l.WithFields(map[string]interface{}{"func": funcName}).Attach(ctx)
			l.Infof("%s done", funcName)
		}(i)
	}
	wg.Wait()
}
