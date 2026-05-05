package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envdiff/internal/loader"
)

func writeTempEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp file %s: %v", p, err)
	}
	return p
}

func TestLoadFiles_Basic(t *testing.T) {
	dir := t.TempDir()
	p := writeTempEnv(t, dir, ".env", "KEY=value\nFOO=bar\n")

	files, err := loader.LoadFiles([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Entries["KEY"] != "value" {
		t.Errorf("expected KEY=value, got %s", files[0].Entries["KEY"])
	}
}

func TestLoadFiles_FileNotFound(t *testing.T) {
	_, err := loader.LoadFiles([]string{"/nonexistent/.env"})
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadFiles_NoPaths(t *testing.T) {
	_, err := loader.LoadFiles([]string{})
	if err == nil {
		t.Fatal("expected error for empty paths, got nil")
	}
}

func TestLoadFiles_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeTempEnv(t, dir, ".env.development", "APP_ENV=development\n")
	p2 := writeTempEnv(t, dir, ".env.production", "APP_ENV=production\n")

	files, err := loader.LoadFiles([]string{p1, p2})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestLoadDir_FindsEnvFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempEnv(t, dir, ".env.staging", "STAGE=true\n")
	writeTempEnv(t, dir, ".env.production", "PROD=true\n")

	files, err := loader.LoadDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) < 2 {
		t.Errorf("expected at least 2 files, got %d", len(files))
	}
}

func TestLoadDir_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	_, err := loader.LoadDir(dir)
	if err == nil {
		t.Fatal("expected error for empty directory, got nil")
	}
}
