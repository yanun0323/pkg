package logs

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"
)

func TestTimerLogger(t *testing.T) {
	{
		writer := &bytes.Buffer{}
		timer := NewTimerLogger(time.Hour, LevelDebug, writer)

		timer.Info("test")
		timer.Info("test")
		timer.Info("test")

		all, err := io.ReadAll(writer)
		if err != nil {
			t.Fatal(err)
		}

		count := strings.Count(string(all), "test")
		if count != 1 {
			t.Fatalf("expected 1 test, got %d", count)
		}
	}

	{
		writer := &bytes.Buffer{}
		timer := NewTimerLogger(time.Second, LevelDebug, writer)

		timer.Info("test")
		timer.Info("test")
		timer.Info("test")

		time.Sleep(time.Second)
		timer.Info("test")

		all, err := io.ReadAll(writer)
		if err != nil {
			t.Fatal(err)
		}

		count := strings.Count(string(all), "test")
		if count != 2 {
			t.Fatalf("expected 2 test, got %d", count)
		}
	}
}
