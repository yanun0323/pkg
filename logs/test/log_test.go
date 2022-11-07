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
	log1.WithField("test", map[string]interface{}{"test": true}).Info("acccess")
	log1.WithField("slack.notify", true).Info("acccess")
	log1.WithField("slack.channelID", 1234234).Info("acccess")
}

func TestLogs_WithEventID(t *testing.T) {
	log := logs.New("logs", 0).WithEventID("fc0d40ae-1801-4bda-9c3b-e1b6956b59b0")

	log.Info("info")
	log.Info("info")
	log.Warn("warn")
	log.Warn("warn")
	log.Error("error")
	log.Error("error")
	log.Fatal("fatal")
	log.Fatal("fatal")
}

func TestLogs_WithPair(t *testing.T) {
	log := logs.New("logs", 4).WithPair("bito", "eth")

	log.Info("info")
	log.Warn("warn")
	log.WithEventID("fc0d40ae-1801-4bda-9c3b-e1b6956b59b0").Error("error")
	log.Fatal("fatal")
}

func TestLogs_WithField(t *testing.T) {
	log := logs.New("logs", 0).WithField("haha", "foo....")

	log.Info("info...")
	log.WithEventID("im eventid").Info("with event id")
}

func TestLogs_WithFields(t *testing.T) {
	log := logs.New("logs", 0).WithFields(map[string]interface{}{
		"foo": 123,
		"bar": "456",
	})

	log.Info("info...")
	log.WithEventID("im eventid").Info("with event id")
}
