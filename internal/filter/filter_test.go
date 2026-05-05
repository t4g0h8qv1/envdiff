package filter_test

import (
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/filter"
)

var sampleResults = []diff.Result{
	{Key: "APP_HOST", Type: diff.Missing, Left: "localhost", Right: ""},
	{Key: "APP_PORT", Type: diff.Conflict, Left: "8080", Right: "9090"},
	{Key: "DB_HOST", Type: diff.Missing, Left: "", Right: "db.prod"},
	{Key: "DB_PASS", Type: diff.Conflict, Left: "secret", Right: "hunter2"},
	{Key: "LOG_LEVEL", Type: diff.Missing, Left: "debug", Right: ""},
}

func TestApply_NoFilter(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{})
	if len(got) != len(sampleResults) {
		t.Errorf("expected %d results, got %d", len(sampleResults), len(got))
	}
}

func TestApply_OnlyMissing(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{OnlyMissing: true})
	for _, r := range got {
		if r.Type != diff.Missing {
			t.Errorf("expected only Missing, got %v for key %s", r.Type, r.Key)
		}
	}
	if len(got) != 3 {
		t.Errorf("expected 3 missing results, got %d", len(got))
	}
}

func TestApply_OnlyConflicts(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{OnlyConflicts: true})
	if len(got) != 2 {
		t.Errorf("expected 2 conflict results, got %d", len(got))
	}
}

func TestApply_KeyPrefix(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{KeyPrefix: "DB_"})
	if len(got) != 2 {
		t.Errorf("expected 2 DB_ results, got %d", len(got))
	}
}

func TestApply_KeyContains(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{KeyContains: "HOST"})
	if len(got) != 2 {
		t.Errorf("expected 2 HOST results, got %d", len(got))
	}
}

func TestApply_CombinedFilters(t *testing.T) {
	got := filter.Apply(sampleResults, filter.Options{
		OnlyMissing: true,
		KeyPrefix:   "APP_",
	})
	if len(got) != 1 {
		t.Errorf("expected 1 result, got %d", len(got))
	}
	if len(got) > 0 && got[0].Key != "APP_HOST" {
		t.Errorf("expected APP_HOST, got %s", got[0].Key)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	got := filter.Apply(nil, filter.Options{OnlyMissing: true})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %d", len(got))
	}
}
