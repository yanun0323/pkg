package logs

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestWithFields(t *testing.T) {
	wg := sync.WaitGroup{}
	l := newWithOutput("test", 0, "stdout")
	wg.Add(100)
	for i := 1; i <= 100; i++ {
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
