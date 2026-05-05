package ignore

import (
	"os"
	"testing"
)

func writeTempIgnore(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "*.envignore")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadFile_ExactKeys(t *testing.T) {
	path := writeTempIgnore(t, "# ignore secrets\nSECRET_KEY\nAPI_TOKEN\n")
	rules, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules.ShouldIgnore("SECRET_KEY") {
		t.Error("expected SECRET_KEY to be ignored")
	}
	if !rules.ShouldIgnore("API_TOKEN") {
		t.Error("expected API_TOKEN to be ignored")
	}
	if rules.ShouldIgnore("DATABASE_URL") {
		t.Error("expected DATABASE_URL not to be ignored")
	}
}

func TestLoadFile_PrefixPatterns(t *testing.T) {
	path := writeTempIgnore(t, "AWS_*\nINTERNAL_*\n")
	rules, err := LoadFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !rules.ShouldIgnore("AWS_ACCESS_KEY") {
		t.Error("expected AWS_ACCESS_KEY to be ignored")
	}
	if !rules.ShouldIgnore("INTERNAL_FLAG") {
		t.Error("expected INTERNAL_FLAG to be ignored")
	}
	if rules.ShouldIgnore("APP_ENV") {
		t.Error("expected APP_ENV not to be ignored")
	}
}

func TestLoadFile_NotFound(t *testing.T) {
	_, err := LoadFile("/nonexistent/.envignore")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestFilterKeys_RemovesIgnored(t *testing.T) {
	rules := NewRules()
	rules.keys["SECRET"] = struct{}{}
	rules.prefixes = []string{"INTERNAL_"}

	env := map[string]string{
		"SECRET":       "abc",
		"INTERNAL_ID":  "42",
		"DATABASE_URL": "postgres://localhost",
	}

	result := rules.FilterKeys(env)

	if _, ok := result["SECRET"]; ok {
		t.Error("SECRET should have been filtered")
	}
	if _, ok := result["INTERNAL_ID"]; ok {
		t.Error("INTERNAL_ID should have been filtered")
	}
	if _, ok := result["DATABASE_URL"]; !ok {
		t.Error("DATABASE_URL should be present")
	}
}

func TestNewRules_Empty(t *testing.T) {
	rules := NewRules()
	if rules.ShouldIgnore("ANYTHING") {
		t.Error("empty rules should not ignore any key")
	}
}
