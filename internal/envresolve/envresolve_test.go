package envresolve_test

import (
	"testing"

	"github.com/user/envdiff/internal/envresolve"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestResolve_NoReferences(t *testing.T) {
	env := makeEnv("HOST", "localhost", "PORT", "8080")
	out, results := envresolve.Resolve(env)
	if out["HOST"] != "localhost" || out["PORT"] != "8080" {
		t.Errorf("unexpected values: %v", out)
	}
	for _, r := range results {
		if r.Expanded {
			t.Errorf("key %s should not be marked expanded", r.Key)
		}
		if r.Err != nil {
			t.Errorf("unexpected error for key %s: %v", r.Key, r.Err)
		}
	}
}

func TestResolve_SimpleReference(t *testing.T) {
	env := makeEnv("SCHEME", "https", "HOST", "example.com", "BASE_URL", "${SCHEME}://${HOST}")
	out, _ := envresolve.Resolve(env)
	if got := out["BASE_URL"]; got != "https://example.com" {
		t.Errorf("BASE_URL = %q, want %q", got, "https://example.com")
	}
}

func TestResolve_NestedReference(t *testing.T) {
	env := makeEnv("A", "hello", "B", "${A}_world", "C", "${B}!")
	out, _ := envresolve.Resolve(env)
	if got := out["C"]; got != "hello_world!" {
		t.Errorf("C = %q, want %q", got, "hello_world!")
	}
}

func TestResolve_UndefinedReference(t *testing.T) {
	env := makeEnv("URL", "${MISSING}/path")
	_, results := envresolve.Resolve(env)
	violations := envresolve.Violations(results)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "URL" {
		t.Errorf("expected violation for URL, got %s", violations[0].Key)
	}
}

func TestResolve_ExpandedFlag(t *testing.T) {
	env := makeEnv("NAME", "world", "GREETING", "hello ${NAME}")
	_, results := envresolve.Resolve(env)
	for _, r := range results {
		if r.Key == "GREETING" && !r.Expanded {
			t.Errorf("GREETING should be marked as expanded")
		}
		if r.Key == "NAME" && r.Expanded {
			t.Errorf("NAME should not be marked as expanded")
		}
	}
}

func TestResolve_OriginalPreserved(t *testing.T) {
	env := makeEnv("HOST", "localhost", "DSN", "db://${HOST}")
	_, results := envresolve.Resolve(env)
	for _, r := range results {
		if r.Key == "DSN" {
			if r.Original != "db://${HOST}" {
				t.Errorf("Original = %q, want %q", r.Original, "db://${HOST}")
			}
			if r.Resolved != "db://localhost" {
				t.Errorf("Resolved = %q, want %q", r.Resolved, "db://localhost")
			}
		}
	}
}

func TestViolations_EmptyWhenNoErrors(t *testing.T) {
	env := makeEnv("A", "1", "B", "2")
	_, results := envresolve.Resolve(env)
	if v := envresolve.Violations(results); len(v) != 0 {
		t.Errorf("expected no violations, got %d", len(v))
	}
}
