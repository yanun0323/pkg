package logs

import (
	"fmt"
	"log/slog"
	"os"
	"sync"
	"testing"
)

func TestVariousLogger(t *testing.T) {
	l := New(LevelDebug, os.Stdout)
	sl := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	t.Log("logger")
	l.WithField("key", "value").Info("logger")
	t.Log("slog")
	sl.With("key", "value").Info("slog")
	t.Log("zap")
}

func TestWithFieldLoop(t *testing.T) {
	wg := sync.WaitGroup{}
	l := New(LevelDebug)
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

func TestWithFieldsLoop(t *testing.T) {
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
