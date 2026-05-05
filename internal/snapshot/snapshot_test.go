package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/snapshot"
)

func sampleResults() []diff.Result {
	return []diff.Result{
		{Key: "DB_HOST", Type: diff.Missing, LeftValue: "localhost", RightValue: ""},
		{Key: "API_KEY", Type: diff.Conflict, LeftValue: "abc", RightValue: "xyz"},
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "snap.json")

	results := sampleResults()
	if err := snapshot.Save(path, "test-label", results); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	s, err := snapshot.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if s.Label != "test-label" {
		t.Errorf("expected label 'test-label', got %q", s.Label)
	}
	if len(s.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(s.Results))
	}
	if s.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if s.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("CreatedAt is in the future")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := snapshot.Load("/nonexistent/path/snap.json")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json{"), 0644)

	_, err := snapshot.Load(path)
	if err == nil {
		t.Error("expected unmarshal error, got nil")
	}
}

func TestCompare_DetectsNewKeys(t *testing.T) {
	base := &snapshot.Snapshot{
		Results: []diff.Result{
			{Key: "DB_HOST", Type: diff.Missing},
		},
	}
	current := []diff.Result{
		{Key: "DB_HOST", Type: diff.Missing},
		{Key: "NEW_KEY", Type: diff.Conflict},
	}

	changed := snapshot.Compare(base, current)
	if len(changed) != 1 || changed[0].Key != "NEW_KEY" {
		t.Errorf("expected only NEW_KEY in changed, got %+v", changed)
	}
}

func TestCompare_DetectsTypeChange(t *testing.T) {
	base := &snapshot.Snapshot{
		Results: []diff.Result{
			{Key: "API_KEY", Type: diff.Missing},
		},
	}
	current := []diff.Result{
		{Key: "API_KEY", Type: diff.Conflict},
	}

	changed := snapshot.Compare(base, current)
	if len(changed) != 1 || changed[0].Key != "API_KEY" {
		t.Errorf("expected API_KEY in changed, got %+v", changed)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	results := sampleResults()
	base := &snapshot.Snapshot{Results: results}

	changed := snapshot.Compare(base, results)
	if len(changed) != 0 {
		t.Errorf("expected no changes, got %+v", changed)
	}
}
