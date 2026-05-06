package lint

import (
	"testing"
)

func makeEnv(file string, kv map[string]string) map[string]map[string]string {
	return map[string]map[string]string{file: kv}
}

func TestCheck_NoViolations(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"DATABASE_URL": "postgres://localhost/db",
		"PORT":         "8080",
	})
	allRules := []Rule{RuleUppercaseKeys, RuleNoSpacesInKeys, RuleNoEmptyValues, RuleKeyFormat}
	violations := Check(envs, allRules)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %d: %+v", len(violations), violations)
	}
}

func TestCheck_LowercaseKey(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"database_url": "postgres://localhost/db",
	})
	violations := Check(envs, []Rule{RuleUppercaseKeys})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != RuleUppercaseKeys {
		t.Errorf("expected rule %s, got %s", RuleUppercaseKeys, violations[0].Rule)
	}
}

func TestCheck_SpaceInKey(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"MY KEY": "value",
	})
	violations := Check(envs, []Rule{RuleNoSpacesInKeys})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != RuleNoSpacesInKeys {
		t.Errorf("expected rule %s, got %s", RuleNoSpacesInKeys, violations[0].Rule)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"SECRET_KEY": "",
	})
	violations := Check(envs, []Rule{RuleNoEmptyValues})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Rule != RuleNoEmptyValues {
		t.Errorf("expected rule %s, got %s", RuleNoEmptyValues, violations[0].Rule)
	}
}

func TestCheck_KeyFormatInvalid(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"1INVALID": "value",
		"VALID_KEY": "ok",
	})
	violations := Check(envs, []Rule{RuleKeyFormat})
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "1INVALID" {
		t.Errorf("expected violation on key '1INVALID', got %q", violations[0].Key)
	}
}

func TestCheck_MultipleFilesMultipleViolations(t *testing.T) {
	envs := map[string]map[string]string{
		".env.dev": {
			"bad key": "value",
			"EMPTY":   "",
		},
		".env.prod": {
			"lowercase": "val",
		},
	}
	allRules := []Rule{RuleUppercaseKeys, RuleNoSpacesInKeys, RuleNoEmptyValues}
	violations := Check(envs, allRules)
	if len(violations) < 3 {
		t.Errorf("expected at least 3 violations, got %d: %+v", len(violations), violations)
	}
}

func TestCheck_OnlyEnabledRulesApply(t *testing.T) {
	envs := makeEnv(".env", map[string]string{
		"lowercase_key": "",
	})
	// Only check empty values — lowercase should not produce a violation
	violations := Check(envs, []Rule{RuleNoEmptyValues})
	for _, v := range violations {
		if v.Rule == RuleUppercaseKeys {
			t.Error("uppercase rule should not have been applied")
		}
	}
}
