package envpromote_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envpromote"
)

func makeEnv(pairs ...string) map[string]string {
	m := map[string]string{}
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestPromote_NewKeysArePromoted(t *testing.T) {
	src := makeEnv("FOO", "bar")
	dst := map[string]string{}

	results, err := envpromote.Promote(src, dst, envpromote.Options{Policy: envpromote.PolicySkip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["FOO"] != "bar" {
		t.Errorf("expected FOO=bar in dst, got %q", dst["FOO"])
	}
	if len(results) != 1 || results[0].Action != "promoted" {
		t.Errorf("expected 1 promoted result, got %+v", results)
	}
}

func TestPromote_PolicySkip_LeavesConflictUnchanged(t *testing.T) {
	src := makeEnv("FOO", "new")
	dst := makeEnv("FOO", "old")

	_, err := envpromote.Promote(src, dst, envpromote.Options{Policy: envpromote.PolicySkip})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["FOO"] != "old" {
		t.Errorf("expected dst FOO to remain 'old', got %q", dst["FOO"])
	}
}

func TestPromote_PolicyOverwrite_ReplacesConflict(t *testing.T) {
	src := makeEnv("FOO", "new")
	dst := makeEnv("FOO", "old")

	results, err := envpromote.Promote(src, dst, envpromote.Options{Policy: envpromote.PolicyOverwrite})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dst["FOO"] != "new" {
		t.Errorf("expected dst FOO='new', got %q", dst["FOO"])
	}
	if results[0].Action != "overwritten" {
		t.Errorf("expected action 'overwritten', got %q", results[0].Action)
	}
}

func TestPromote_PolicyError_ReturnsErrorOnConflict(t *testing.T) {
	src := makeEnv("FOO", "new")
	dst := makeEnv("FOO", "old")

	_, err := envpromote.Promote(src, dst, envpromote.Options{Policy: envpromote.PolicyError})
	if err == nil {
		t.Fatal("expected error on conflict, got nil")
	}
}

func TestPromote_AllowKeys_FiltersOtherKeys(t *testing.T) {
	src := makeEnv("FOO", "1", "BAR", "2")
	dst := map[string]string{}

	_, err := envpromote.Promote(src, dst, envpromote.Options{
		AllowKeys: []string{"FOO"},
		Policy:    envpromote.PolicySkip,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["BAR"]; ok {
		t.Error("BAR should not have been promoted")
	}
	if dst["FOO"] != "1" {
		t.Errorf("expected FOO=1, got %q", dst["FOO"])
	}
}

func TestPromote_DenyKeys_ExcludesKeys(t *testing.T) {
	src := makeEnv("SECRET", "s3cr3t", "HOST", "localhost")
	dst := map[string]string{}

	_, err := envpromote.Promote(src, dst, envpromote.Options{
		DenyKeys: []string{"SECRET"},
		Policy:   envpromote.PolicySkip,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := dst["SECRET"]; ok {
		t.Error("SECRET should have been denied")
	}
	if dst["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", dst["HOST"])
	}
}

func TestSummary_ReturnsReadableString(t *testing.T) {
	results := []envpromote.Result{
		{Key: "A", Action: "promoted"},
		{Key: "B", Action: "promoted"},
		{Key: "C", Action: "overwritten"},
		{Key: "D", Action: "skipped"},
	}
	got := envpromote.Summary(results)
	expected := "2 promoted, 1 overwritten, 1 skipped"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestSummary_NothingToPromote(t *testing.T) {
	got := envpromote.Summary(nil)
	if got != "nothing to promote" {
		t.Errorf("expected 'nothing to promote', got %q", got)
	}
}
