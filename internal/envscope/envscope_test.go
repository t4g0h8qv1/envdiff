package envscope_test

import (
	"testing"

	"github.com/user/envdiff/internal/envscope"
)

func makeScope(name string, pairs ...string) envscope.Scope {
	vars := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		vars[pairs[i]] = pairs[i+1]
	}
	return envscope.Scope{Name: name, Vars: vars}
}

func TestResolve_FirstScopeWins(t *testing.T) {
	r := envscope.New(
		makeScope("prod", "DB_HOST", "prod-db"),
		makeScope("dev", "DB_HOST", "localhost"),
	)
	val, scope, err := r.Resolve("DB_HOST")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "prod-db" {
		t.Errorf("expected prod-db, got %q", val)
	}
	if scope != "prod" {
		t.Errorf("expected scope prod, got %q", scope)
	}
}

func TestResolve_FallsBackToLowerScope(t *testing.T) {
	r := envscope.New(
		makeScope("prod", "APP_PORT", "8080"),
		makeScope("dev", "APP_PORT", "3000", "DEBUG", "true"),
	)
	val, scope, err := r.Resolve("DEBUG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "true" {
		t.Errorf("expected true, got %q", val)
	}
	if scope != "dev" {
		t.Errorf("expected scope dev, got %q", scope)
	}
}

func TestResolve_KeyNotFound(t *testing.T) {
	r := envscope.New(makeScope("prod", "FOO", "bar"))
	_, _, err := r.Resolve("MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestResolveAll_MergesScopes(t *testing.T) {
	r := envscope.New(
		makeScope("prod", "HOST", "prod-host", "PORT", "443"),
		makeScope("dev", "HOST", "localhost", "DEBUG", "true"),
	)
	all := r.ResolveAll()
	if all["HOST"] != "prod-host" {
		t.Errorf("expected prod-host, got %q", all["HOST"])
	}
	if all["DEBUG"] != "true" {
		t.Errorf("expected true, got %q", all["DEBUG"])
	}
	if all["PORT"] != "443" {
		t.Errorf("expected 443, got %q", all["PORT"])
	}
}

func TestFindConflicts_DetectsConflict(t *testing.T) {
	r := envscope.New(
		makeScope("prod", "DB_URL", "postgres://prod"),
		makeScope("dev", "DB_URL", "postgres://dev"),
	)
	conflicts := r.FindConflicts()
	if len(conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(conflicts))
	}
	if conflicts[0].Key != "DB_URL" {
		t.Errorf("expected conflict on DB_URL, got %q", conflicts[0].Key)
	}
}

func TestFindConflicts_NoConflictWhenSameValue(t *testing.T) {
	r := envscope.New(
		makeScope("prod", "LOG_LEVEL", "info"),
		makeScope("dev", "LOG_LEVEL", "info"),
	)
	conflicts := r.FindConflicts()
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}

func TestFindConflicts_SingleScopeNoConflict(t *testing.T) {
	r := envscope.New(makeScope("prod", "KEY", "value"))
	conflicts := r.FindConflicts()
	if len(conflicts) != 0 {
		t.Errorf("expected no conflicts, got %d", len(conflicts))
	}
}
