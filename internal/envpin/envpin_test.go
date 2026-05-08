package envpin_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/envpin"
)

func writeTempPin(t *testing.T, pf envpin.PinFile) string {
	t.Helper()
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(t.TempDir(), "pins.json")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSaveAndLoad(t *testing.T) {
	env := map[string]string{"APP_ENV": "production", "LOG_LEVEL": "info"}
	path := filepath.Join(t.TempDir(), "pins.json")

	if err := envpin.Save(path, env, []string{"APP_ENV", "LOG_LEVEL"}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	pf, err := envpin.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(pf.Pins) != 2 {
		t.Fatalf("expected 2 pins, got %d", len(pf.Pins))
	}
}

func TestSave_MissingKey(t *testing.T) {
	env := map[string]string{"APP_ENV": "production"}
	path := filepath.Join(t.TempDir(), "pins.json")
	err := envpin.Save(path, env, []string{"MISSING_KEY"})
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestCheck_NoViolations(t *testing.T) {
	pf := envpin.PinFile{
		Pins: []envpin.Pin{
			{Key: "APP_ENV", Expected: "production"},
		},
	}
	env := map[string]string{"APP_ENV": "production"}
	violations := envpin.Check(pf, env)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestCheck_ValueChanged(t *testing.T) {
	pf := envpin.PinFile{
		Pins: []envpin.Pin{
			{Key: "APP_ENV", Expected: "production"},
		},
	}
	env := map[string]string{"APP_ENV": "staging"}
	violations := envpin.Check(pf, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Actual != "staging" {
		t.Errorf("unexpected actual value: %s", violations[0].Actual)
	}
}

func TestCheck_MissingKey(t *testing.T) {
	pf := envpin.PinFile{
		Pins: []envpin.Pin{
			{Key: "DB_HOST", Expected: "localhost"},
		},
	}
	env := map[string]string{}
	violations := envpin.Check(pf, env)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if !violations[0].Missing {
		t.Error("expected Missing=true")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := envpin.Load("/nonexistent/pins.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
