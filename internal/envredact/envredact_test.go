package envredact_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/envredact"
)

func TestIsSensitive_MatchesDefaultPatterns(t *testing.T) {
	r := envredact.New(nil)
	sensitive := []string{"DB_PASSWORD", "API_KEY", "AUTH_TOKEN", "SECRET_KEY", "AWS_ACCESS_KEY"}
	for _, key := range sensitive {
		if !r.IsSensitive(key) {
			t.Errorf("expected %q to be sensitive", key)
		}
	}
}

func TestIsSensitive_SafeKeys(t *testing.T) {
	r := envredact.New(nil)
	safe := []string{"APP_ENV", "PORT", "LOG_LEVEL", "BASE_URL"}
	for _, key := range safe {
		if r.IsSensitive(key) {
			t.Errorf("expected %q to NOT be sensitive", key)
		}
	}
}

func TestRedact_ReplacesOnlySensitiveValues(t *testing.T) {
	r := envredact.New(nil)
	env := map[string]string{
		"APP_ENV":     "production",
		"DB_PASSWORD": "supersecret",
		"PORT":        "8080",
		"API_KEY":     "abc123",
	}

	result := r.Redact(env)

	if result["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be redacted, got %q", result["APP_ENV"])
	}
	if result["PORT"] != "8080" {
		t.Errorf("PORT should not be redacted, got %q", result["PORT"])
	}
	if result["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD should be redacted, got %q", result["DB_PASSWORD"])
	}
	if result["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY should be redacted, got %q", result["API_KEY"])
	}
}

func TestRedact_OriginalMapUnmodified(t *testing.T) {
	r := envredact.New(nil)
	env := map[string]string{"API_KEY": "original-value"}
	_ = r.Redact(env)
	if env["API_KEY"] != "original-value" {
		t.Error("original map should not be modified")
	}
}

func TestRedactAll_AppliesAcrossMultipleMaps(t *testing.T) {
	r := envredact.New(nil)
	envs := map[string]map[string]string{
		"staging": {"APP_ENV": "staging", "DB_PASSWORD": "stg-pass"},
		"prod":    {"APP_ENV": "production", "DB_PASSWORD": "prd-pass"},
	}

	result := r.RedactAll(envs)

	for env, m := range result {
		if m["DB_PASSWORD"] != "[REDACTED]" {
			t.Errorf("[%s] DB_PASSWORD should be redacted", env)
		}
		if m["APP_ENV"] == "[REDACTED]" {
			t.Errorf("[%s] APP_ENV should not be redacted", env)
		}
	}
}

func TestNew_CustomPatterns(t *testing.T) {
	r := envredact.New([]string{"INTERNAL"})
	if !r.IsSensitive("MY_INTERNAL_FLAG") {
		t.Error("expected MY_INTERNAL_FLAG to be sensitive with custom pattern")
	}
	if r.IsSensitive("API_KEY") {
		t.Error("API_KEY should not be sensitive with custom-only pattern set")
	}
}
