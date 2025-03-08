package logs

import "testing"

func TestFileWriter(t *testing.T) {
	writer := FileWriter(".", "test")
	writer.Write([]byte("test"))
	if err := writer.Sync(); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	if err := writer.Remove(); err != nil {
		t.Fatalf("remove failed: %v", err)
	}
}
