package envcast

import (
	"testing"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestCastMap_AllValid(t *testing.T) {
	env := makeEnv("PORT", "8080", "DEBUG", "true", "RATIO", "1.5", "NAME", "app")
	hints := map[string]string{"PORT": "int", "DEBUG": "bool", "RATIO": "float", "NAME": "string"}
	results := CastMap(env, hints)
	if v := Violations(results); len(v) != 0 {
		t.Errorf("expected no violations, got %d: %v", len(v), v[0].Error)
	}
}

func TestCastMap_InvalidInt(t *testing.T) {
	env := makeEnv("PORT", "not-a-number")
	hints := map[string]string{"PORT": "int"}
	results := CastMap(env, hints)
	if v := Violations(results); len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestCastMap_InvalidBool(t *testing.T) {
	env := makeEnv("FLAG", "yes_please")
	hints := map[string]string{"FLAG": "bool"}
	results := CastMap(env, hints)
	if v := Violations(results); len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
}

func TestCastMap_MissingKey(t *testing.T) {
	env := makeEnv()
	hints := map[string]string{"MISSING": "string"}
	results := CastMap(env, hints)
	if v := Violations(results); len(v) != 1 {
		t.Errorf("expected 1 violation for missing key, got %d", len(v))
	}
}

func TestCastMap_UnknownTypeHint(t *testing.T) {
	env := makeEnv("KEY", "value")
	hints := map[string]string{"KEY": "uuid"}
	results := CastMap(env, hints)
	if v := Violations(results); len(v) != 1 {
		t.Errorf("expected 1 violation for unknown type hint, got %d", len(v))
	}
}

func TestViolations_FiltersCorrectly(t *testing.T) {
	env := makeEnv("PORT", "8080", "DEBUG", "bad")
	hints := map[string]string{"PORT": "int", "DEBUG": "bool"}
	results := CastMap(env, hints)
	v := Violations(results)
	if len(v) != 1 {
		t.Errorf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DEBUG" {
		t.Errorf("expected violation for DEBUG, got %q", v[0].Key)
	}
}
