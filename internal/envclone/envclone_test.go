package envclone_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/envclone"
)

func tempTarget(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), ".env.clone")
}

func TestClone_AllKeys(t *testing.T) {
	src := map[string]string{"APP_NAME": "myapp", "PORT": "8080"}
	target := tempTarget(t)

	res, err := envclone.Clone(src, target, envclone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 2 {
		t.Errorf("expected 2 cloned, got %d", res.Cloned)
	}
	data, _ := os.ReadFile(target)
	if !strings.Contains(string(data), "APP_NAME=myapp") {
		t.Errorf("expected APP_NAME in output, got: %s", data)
	}
}

func TestClone_PrefixFilter(t *testing.T) {
	src := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "APP_ENV": "prod"}
	target := tempTarget(t)

	res, err := envclone.Clone(src, target, envclone.Options{KeyPrefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 2 || res.Skipped != 1 {
		t.Errorf("expected 2 cloned, 1 skipped; got %d cloned, %d skipped", res.Cloned, res.Skipped)
	}
	data, _ := os.ReadFile(target)
	if strings.Contains(string(data), "APP_ENV") {
		t.Errorf("APP_ENV should have been filtered out")
	}
}

func TestClone_Redact(t *testing.T) {
	src := map[string]string{"SECRET_KEY": "supersecret"}
	target := tempTarget(t)

	_, err := envclone.Clone(src, target, envclone.Options{Redact: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(target)
	if strings.Contains(string(data), "supersecret") {
		t.Errorf("expected value to be redacted")
	}
	if !strings.Contains(string(data), "SECRET_KEY=") {
		t.Errorf("expected key to still be present")
	}
}

func TestClone_OverwriteProtection(t *testing.T) {
	src := map[string]string{"KEY": "val"}
	target := tempTarget(t)
	os.WriteFile(target, []byte("existing"), 0o644)

	_, err := envclone.Clone(src, target, envclone.Options{})
	if err == nil {
		t.Fatal("expected error when target exists and OverwriteExisting is false")
	}
}

func TestClone_OverwriteAllowed(t *testing.T) {
	src := map[string]string{"KEY": "newval"}
	target := tempTarget(t)
	os.WriteFile(target, []byte("old=data\n"), 0o644)

	_, err := envclone.Clone(src, target, envclone.Options{OverwriteExisting: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(target)
	if !strings.Contains(string(data), "KEY=newval") {
		t.Errorf("expected overwritten content, got: %s", data)
	}
}

func TestClone_EmptySrc(t *testing.T) {
	target := tempTarget(t)
	res, err := envclone.Clone(map[string]string{}, target, envclone.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Cloned != 0 {
		t.Errorf("expected 0 cloned, got %d", res.Cloned)
	}
}
