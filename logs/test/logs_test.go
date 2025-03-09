package test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"testing"

	"github.com/yanun0323/pkg/logs"
)

func TestGet(t *testing.T) {
	log := logs.Get(context.Background())
	log.Info("Test")
}

func TestLogOutput(t *testing.T) {
	w2 := logs.FileWriter(".", "dir_dot")
	w3 := logs.FileWriter("", "dir_empty")
	w4 := logs.FileWriter("hello", "dir_wrong")

	log1 := logs.New(logs.LevelInfo, os.Stdout)
	log2 := logs.New(logs.LevelInfo, os.Stdout, w2)
	log3 := logs.New(logs.LevelInfo, os.Stdout, w3)
	log4 := logs.New(logs.LevelInfo, os.Stdout, w4)

	t.Logf("log1 = %p, log2 = %p, log3 = %p, log4 = %p", log1, log2, log3, log4)
	log1.Info("info")
	log2.Info("info")
	log3.Info("info")
	log4.Info("info")

	if err := w2.Remove(); err != nil {
		t.Errorf("remove w2 failed: %v", err)
	}

	if err := w3.Remove(); err != nil {
		t.Errorf("remove w3 failed: %v", err)
	}

	if err := w4.Remove(); err != nil {
		t.Errorf("remove w4 failed: %v", err)
	}
}

func TestLogs(t *testing.T) {
	log1 := logs.New(logs.LevelInfo)
	log2 := logs.New(logs.LevelInfo)

	t.Logf("log1 = %p, log2 = %p", log1, log2)
	log1.Info("info")
	log2.Info("info")
}

func TestMap(t *testing.T) {
	log1 := logs.New(logs.LevelInfo)
	log1.WithField("test", map[string]interface{}{"test": true}).Info("access")
}

func TestLogs_WithField(t *testing.T) {
	log := logs.New(logs.LevelInfo).WithField("hello", "foo....")

	log.Info("hello field info...")
	log.WithField("user_id", "i'm user").WithField("info_id", "i'm order").Info("with user id")
}

func TestLogs_WithFields(t *testing.T) {
	log := logs.New(logs.LevelInfo).WithFields(map[string]interface{}{
		"foo": 123,
		"bar": "456",
	})

	log.Info("info...")
	log.WithField("user_id", "i'm user").Info("with user id")
}

func TestLogs_Fatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		logs.New(logs.LevelInfo).Fatal("fatal")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestLogs_Fatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestLogs_Level(t *testing.T) {
	a := NewAssert(t)
	writer := NewMockWriter()
	log := logs.New(logs.LevelWarn, writer)
	log.Info("[LEVEL] info")
	a.True(len(writer.ReadAndClean()) == 0)

	log.WithField("user_id", "i'm user").
		Warn("[LEVEL] with user id")

	a.True(len(writer.ReadAndClean()) != 0)
}

type MockWriter struct {
	buf bytes.Buffer
}

func NewMockWriter() *MockWriter {
	return &MockWriter{}
}

func (w *MockWriter) Write(p []byte) (n int, err error) {
	_, _ = os.Stdout.Write(p)
	return w.buf.Write(p)
}

func (w *MockWriter) String() string {
	return w.buf.String()
}

func (w *MockWriter) Remove() error {
	return nil
}

func (w *MockWriter) Sync() error {
	return nil
}

func (w *MockWriter) ReadAndClean() string {
	s := w.String()
	w.buf.Reset()
	return s
}

func (w *MockWriter) Clean() {
	w.buf.Reset()
}
