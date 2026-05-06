package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/audit"
	"github.com/user/envdiff/internal/diff"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Status: "missing", Left: "localhost", Right: ""},
		{Key: "API_KEY", Status: "conflict", Left: "abc", Right: "xyz"},
		{Key: "PORT", Status: "ok", Left: "8080", Right: "8080"},
	}
}

func TestNewEntry_Summary(t *testing.T) {
	files := []string{".env.dev", ".env.prod"}
	results := sampleResults()
	entry := audit.NewEntry(files, results)

	if entry.Summary.Total != 3 {
		t.Errorf("expected Total=3, got %d", entry.Summary.Total)
	}
	if entry.Summary.Missing != 1 {
		t.Errorf("expected Missing=1, got %d", entry.Summary.Missing)
	}
	if entry.Summary.Conflicts != 1 {
		t.Errorf("expected Conflicts=1, got %d", entry.Summary.Conflicts)
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
	if len(entry.Files) != 2 {
		t.Errorf("expected 2 files, got %d", len(entry.Files))
	}
}

func TestAppendAndReadAll(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	files := []string{".env.staging"}
	results := sampleResults()

	e1 := audit.NewEntry(files, results[:1])
	e2 := audit.NewEntry(files, results[1:2])

	if err := audit.Append(logPath, e1); err != nil {
		t.Fatalf("Append e1: %v", err)
	}
	if err := audit.Append(logPath, e2); err != nil {
		t.Fatalf("Append e2: %v", err)
	}

	entries, err := audit.ReadAll(logPath)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Summary.Missing != 1 {
		t.Errorf("entry[0] Missing: expected 1, got %d", entries[0].Summary.Missing)
	}
	if entries[1].Summary.Conflicts != 1 {
		t.Errorf("entry[1] Conflicts: expected 1, got %d", entries[1].Summary.Conflicts)
	}
}

func TestReadAll_FileNotFound(t *testing.T) {
	_, err := audit.ReadAll("/nonexistent/audit.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestAppend_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "new_audit.log")

	entry := audit.NewEntry([]string{".env"}, nil)
	if err := audit.Append(logPath, entry); err != nil {
		t.Fatalf("Append: %v", err)
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("expected log file to be created")
	}
}

func TestNewEntry_TimestampIsUTC(t *testing.T) {
	before := time.Now().UTC()
	entry := audit.NewEntry(nil, nil)
	after := time.Now().UTC()

	if entry.Timestamp.Before(before) || entry.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", entry.Timestamp, before, after)
	}
}
