package envcheck

import (
	"strings"
	"testing"
)

func makeRef(keys ...string) map[string]string {
	m := make(map[string]string, len(keys))
	for _, k := range keys {
		m[k] = "ref_value"
	}
	return m
}

func TestCheck_NoViolations(t *testing.T) {
	ref := makeRef("DB_HOST", "DB_PORT")
	targets := map[string]map[string]string{
		"staging": {"DB_HOST": "localhost", "DB_PORT": "5432"},
	}
	results := Check(ref, targets, Options{})
	if len(results) != 0 {
		t.Fatalf("expected no results, got %d", len(results))
	}
}

func TestCheck_MissingKey(t *testing.T) {
	ref := makeRef("DB_HOST", "DB_PORT", "DB_NAME")
	targets := map[string]map[string]string{
		"prod": {"DB_HOST": "prod-host", "DB_PORT": "5432"},
	}
	results := Check(ref, targets, Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_NAME" || !results[0].Missing {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestCheck_EmptyValueNotFlaggedByDefault(t *testing.T) {
	ref := makeRef("API_KEY")
	targets := map[string]map[string]string{
		"dev": {"API_KEY": ""},
	}
	results := Check(ref, targets, Options{RequireNonEmpty: false})
	if len(results) != 0 {
		t.Fatalf("expected no results without RequireNonEmpty, got %d", len(results))
	}
}

func TestCheck_EmptyValueFlaggedWhenRequired(t *testing.T) {
	ref := makeRef("API_KEY")
	targets := map[string]map[string]string{
		"dev": {"API_KEY": ""},
	}
	results := Check(ref, targets, Options{RequireNonEmpty: true})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Empty {
		t.Errorf("expected Empty=true, got %+v", results[0])
	}
}

func TestCheck_MultipleTargets(t *testing.T) {
	ref := makeRef("SECRET")
	targets := map[string]map[string]string{
		"staging": {"SECRET": "abc"},
		"prod":    {},
	}
	results := Check(ref, targets, Options{})
	if len(results) != 1 {
		t.Fatalf("expected 1 missing result, got %d", len(results))
	}
	if results[0].EnvName != "prod" || !results[0].Missing {
		t.Errorf("unexpected result: %+v", results[0])
	}
}

func TestFormat_NoResults(t *testing.T) {
	out := Format(nil)
	if out != "All keys present and valid." {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestFormat_WithResults(t *testing.T) {
	results := []Result{
		{Key: "DB_HOST", EnvName: "prod", Missing: true},
		{Key: "API_KEY", EnvName: "staging", Empty: true},
	}
	out := Format(results)
	if !strings.Contains(out, "[MISSING]") {
		t.Error("expected MISSING tag in output")
	}
	if !strings.Contains(out, "[EMPTY]") {
		t.Error("expected EMPTY tag in output")
	}
}
