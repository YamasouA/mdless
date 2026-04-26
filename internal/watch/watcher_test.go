package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWatchReturnsEventOnFileWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte("before"), 0o644); err != nil {
		t.Fatal(err)
	}

	result := make(chan any, 1)
	go func() {
		result <- Watch(path)()
	}()

	time.Sleep(50 * time.Millisecond)
	if err := os.WriteFile(path, []byte("after"), 0o644); err != nil {
		t.Fatal(err)
	}

	select {
	case msg := <-result:
		event, ok := msg.(Event)
		if !ok {
			t.Fatalf("msg = %#v, want Event", msg)
		}
		if event.Path != path {
			t.Fatalf("Path = %q, want %q", event.Path, path)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for watch event")
	}
}
