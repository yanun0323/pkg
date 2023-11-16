package test

import (
	"context"
	"testing"

	"github.com/yanun0323/pkg/logs"
)

func TestGet(t *testing.T) {
	log := logs.Get(context.Background())
	log.Info("Test")
}

func TestLogs(t *testing.T) {
	log1 := logs.New("logs", 0)
	log2 := logs.New("wallet-api", 0)

	t.Logf("log1 = %p, log2 = %p", log1, log2)
	log1.Info("info")
	log2.Info("info")
}

func TestMap(t *testing.T) {
	log1 := logs.New("logs", 0)
	log1.WithField("test", map[string]interface{}{"test": true}).Info("access")
}

func TestLogs_WithFunc(t *testing.T) {
	log := logs.New("logs", 0).WithFunc("WithFunc")

	log.Info("info")
	log.Warn("warn")
	log.Error("error")
}

func TestLogs_Fatal(t *testing.T) {
	log := logs.New("logs", 0)
	t.Cleanup(func() {
		log.Fatal("fatal")
	})
}

func TestLogs_WithField(t *testing.T) {
	log := logs.New("logs", 0).WithField("hello", "foo....")

	log.Info("info...")
	log.WithField("user_id", "im user").Info("with user id")
}

func TestLogs_WithFields(t *testing.T) {
	log := logs.New("logs", 0).WithFields(map[string]interface{}{
		"foo": 123,
		"bar": "456",
	})

	log.Info("info...")
	log.WithField("user_id", "im user").Info("with user id")
}
