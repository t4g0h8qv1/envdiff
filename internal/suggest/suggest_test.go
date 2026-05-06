package suggest_test

import (
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/suggest"
)

func TestGenerate_Empty(t *testing.T) {
	results := []diff.Result{}
	suggestions := suggest.Generate(results)
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
}

func TestGenerate_MissingKey(t *testing.T) {
	results := []diff.Result{
		{Key: "DB_HOST", Kind: "missing", File: ".env.production"},
	}
	suggestions := suggest.Generate(results)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	s := suggestions[0]
	if s.Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", s.Key)
	}
	if s.Kind != "missing" {
		t.Errorf("expected kind missing, got %s", s.Kind)
	}
	if !strings.Contains(s.Action, "DB_HOST") {
		t.Errorf("action should mention key, got: %s", s.Action)
	}
}

func TestGenerate_Conflict(t *testing.T) {
	results := []diff.Result{
		{Key: "API_URL", Kind: "conflict", Files: []string{".env.staging", ".env.production"}},
	}
	suggestions := suggest.Generate(results)
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	s := suggestions[0]
	if s.Kind != "conflict" {
		t.Errorf("expected kind conflict, got %s", s.Kind)
	}
	if !strings.Contains(s.Action, "API_URL") {
		t.Errorf("action should mention key, got: %s", s.Action)
	}
}

func TestFormat_NoSuggestions(t *testing.T) {
	out := suggest.Format(nil)
	if !strings.Contains(out, "consistent") {
		t.Errorf("expected consistent message, got: %s", out)
	}
}

func TestFormat_WithSuggestions(t *testing.T) {
	suggestions := []suggest.Suggestion{
		{Key: "SECRET", Kind: "missing", File: ".env.prod", Action: "Add key \"SECRET\" to .env.prod"},
		{Key: "PORT", Kind: "conflict", File: ".env.dev, .env.prod", Action: "Resolve conflicting values for key \"PORT\""},
	}
	out := suggest.Format(suggestions)
	if !strings.Contains(out, "Suggestions:") {
		t.Errorf("expected header, got: %s", out)
	}
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING label, got: %s", out)
	}
	if !strings.Contains(out, "CONFLICT") {
		t.Errorf("expected CONFLICT label, got: %s", out)
	}
}
