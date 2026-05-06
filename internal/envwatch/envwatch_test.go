package envwatch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/envwatch"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeTempEnv: %v", err)
	}
	return p
}

func TestWatcher_DetectsChange(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "KEY=old\n")

	w := envwatch.New([]string{f}, 50*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(80 * time.Millisecond)

	if err := os.WriteFile(f, []byte("KEY=new\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	select {
	case ev := <-w.Events:
		if ev.File != f {
			t.Errorf("expected file %s, got %s", f, ev.File)
		}
		if len(ev.Results) == 0 {
			t.Error("expected at least one diff result")
		}
	case <-time.After(500 * time.Millisecond):
		t.Error("timed out waiting for change event")
	}
}

func TestWatcher_NoEventWhenUnchanged(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "KEY=same\n")

	w := envwatch.New([]string{f}, 50*time.Millisecond)
	w.Start()
	defer w.Stop()

	time.Sleep(200 * time.Millisecond)

	select {
	case ev := <-w.Events:
		t.Errorf("unexpected event: %+v", ev)
	default:
		// pass
	}
}

func TestWatcher_StopStopsPolling(t *testing.T) {
	dir := t.TempDir()
	f := writeTempEnv(t, dir, ".env", "A=1\n")

	w := envwatch.New([]string{f}, 30*time.Millisecond)
	w.Start()
	w.Stop()

	time.Sleep(100 * time.Millisecond)
	if err := os.WriteFile(f, []byte("A=2\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	time.Sleep(100 * time.Millisecond)

	select {
	case ev := <-w.Events:
		t.Errorf("received event after stop: %+v", ev)
	default:
		// pass
	}
}
