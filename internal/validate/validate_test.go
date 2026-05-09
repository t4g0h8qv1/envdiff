package validate

import (
	"testing"
)

func makeEnvMaps(data map[string]map[string]string) map[string]map[string]string {
	return data
}

func TestCheck_NoViolations(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.production": {"PORT": "8080", "APP_ENV": "production"},
	})
	rules := []Rule{
		{Key: "PORT", Required: true, Pattern: `^\d+$`},
		{Key: "APP_ENV", Required: true},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestCheck_MissingRequiredKey(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.staging": {"APP_ENV": "staging"},
	})
	rules := []Rule{
		{Key: "PORT", Required: true},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %s", violations[0].Key)
	}
	if violations[0].Message != "required key is missing" {
		t.Errorf("unexpected message: %s", violations[0].Message)
	}
}

func TestCheck_PatternMismatch(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.production": {"PORT": "not-a-number"},
	})
	rules := []Rule{
		{Key: "PORT", Pattern: `^\d+$`},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "PORT" {
		t.Errorf("expected violation for PORT, got %s", violations[0].Key)
	}
}

func TestCheck_MultipleFilesMultipleViolations(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.staging":    {"APP_ENV": "staging"},
		".env.production": {"APP_ENV": "production"},
	})
	rules := []Rule{
		{Key: "PORT", Required: true},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 2 {
		t.Errorf("expected 2 violations (one per file), got %d", len(violations))
	}
}

func TestCheck_OptionalKeySkippedWhenAbsent(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.local": {"APP_ENV": "local"},
	})
	rules := []Rule{
		{Key: "OPTIONAL_KEY", Required: false, Pattern: `^[a-z]+$`},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for absent optional key, got %d", len(violations))
	}
}

func TestCheck_ViolationContainsFileName(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.production": {"APP_ENV": "production"},
	})
	rules := []Rule{
		{Key: "PORT", Required: true},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].File != ".env.production" {
		t.Errorf("expected violation file to be '.env.production', got %q", violations[0].File)
	}
}

func TestCheck_EmptyEnvMaps(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{})
	rules := []Rule{
		{Key: "PORT", Required: true},
	}
	violations := Check(envMaps, rules)
	if len(violations) != 0 {
		t.Errorf("expected no violations for empty env maps, got %d", len(violations))
	}
}

func TestCheck_EmptyRules(t *testing.T) {
	envMaps := makeEnvMaps(map[string]map[string]string{
		".env.production": {"PORT": "8080"},
	})
	violations := Check(envMaps, []Rule{})
	if len(violations) != 0 {
		t.Errorf("expected no violations for empty rules, got %d", len(violations))
	}
}
