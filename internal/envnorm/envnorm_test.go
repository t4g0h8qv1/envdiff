package envnorm_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envnorm"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestNormalize_NoChanges(t *testing.T) {
	env := makeEnv("DB_HOST", "localhost", "PORT", "5432")
	opts := envnorm.DefaultOptions()
	result, violations := envnorm.Normalize(env, opts)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
	if result["DB_HOST"] != "localhost" {
		t.Errorf("unexpected value for DB_HOST: %s", result["DB_HOST"])
	}
}

func TestNormalize_UppercasesKeys(t *testing.T) {
	env := makeEnv("db_host", "localhost", "api_key", "secret")
	opts := envnorm.DefaultOptions()
	result, violations := envnorm.Normalize(env, opts)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
	if _, ok := result["DB_HOST"]; !ok {
		t.Error("expected DB_HOST in result")
	}
	if _, ok := result["db_host"]; ok {
		t.Error("expected db_host to be removed")
	}
}

func TestNormalize_TrimsValues(t *testing.T) {
	env := makeEnv("HOST", "  localhost  ", "PORT", "\t8080\t")
	opts := envnorm.DefaultOptions()
	result, violations := envnorm.Normalize(env, opts)
	if len(violations) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(violations))
	}
	if result["HOST"] != "localhost" {
		t.Errorf("expected trimmed value, got %q", result["HOST"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("expected trimmed value, got %q", result["PORT"])
	}
}

func TestNormalize_RemovesEmptyValues(t *testing.T) {
	env := makeEnv("EMPTY", "", "PRESENT", "yes")
	opts := envnorm.DefaultOptions()
	opts.RemoveEmpty = true
	result, violations := envnorm.Normalize(env, opts)
	if _, ok := result["EMPTY"]; ok {
		t.Error("expected EMPTY key to be removed")
	}
	if result["PRESENT"] != "yes" {
		t.Errorf("expected PRESENT to remain")
	}
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Reason != "empty value removed" {
		t.Errorf("unexpected reason: %s", violations[0].Reason)
	}
}

func TestNormalize_QuotesValues(t *testing.T) {
	env := makeEnv("MSG", "hello world")
	opts := envnorm.DefaultOptions()
	opts.QuoteValues = true
	result, violations := envnorm.Normalize(env, opts)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if result["MSG"] != `"hello world"` {
		t.Errorf("unexpected quoted value: %s", result["MSG"])
	}
}

func TestNormalize_AlreadyQuotedSkipped(t *testing.T) {
	env := makeEnv("MSG", `"already quoted"`)
	opts := envnorm.DefaultOptions()
	opts.QuoteValues = true
	_, violations := envnorm.Normalize(env, opts)
	for _, v := range violations {
		if v.Reason == "value quoted" {
			t.Error("should not re-quote an already quoted value")
		}
	}
}
