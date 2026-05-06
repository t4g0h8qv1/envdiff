package envtemplate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourorg/envdiff/internal/envtemplate"
)

func TestGenerate_Empty(t *testing.T) {
	lines := envtemplate.Generate(nil, envtemplate.Options{})
	if len(lines) != 0 {
		t.Fatalf("expected no lines, got %d", len(lines))
	}
}

func TestGenerate_SingleMap(t *testing.T) {
	m := map[string]string{"DB_HOST": "localhost", "APP_PORT": "8080"}
	lines := envtemplate.Generate([]map[string]string{m}, envtemplate.Options{})
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "APP_PORT=" {
		t.Errorf("expected APP_PORT= first (sorted), got %q", lines[0])
	}
	if lines[1] != "DB_HOST=" {
		t.Errorf("expected DB_HOST= second, got %q", lines[1])
	}
}

func TestGenerate_Placeholder(t *testing.T) {
	m := map[string]string{"SECRET": "abc"}
	lines := envtemplate.Generate([]map[string]string{m}, envtemplate.Options{Placeholder: "CHANGE_ME"})
	if lines[0] != "SECRET=CHANGE_ME" {
		t.Errorf("unexpected line: %q", lines[0])
	}
}

func TestGenerate_MergesMultipleMaps(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "99", "C": "3"}
	lines := envtemplate.Generate([]map[string]string{a, b}, envtemplate.Options{})
	if len(lines) != 3 {
		t.Fatalf("expected 3 unique keys, got %d", len(lines))
	}
}

func TestGenerate_IncludeComments(t *testing.T) {
	m := map[string]string{"KEY": "val"}
	lines := envtemplate.Generate([]map[string]string{m}, envtemplate.Options{IncludeComments: true})
	if !strings.HasPrefix(lines[0], "#") {
		t.Errorf("expected first line to be a comment, got %q", lines[0])
	}
}

func TestWrite_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, ".env.template")
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := envtemplate.Write(out, []map[string]string{m}, envtemplate.Options{}); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("could not read output file: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "BAZ=") || !strings.Contains(content, "FOO=") {
		t.Errorf("unexpected file content:\n%s", content)
	}
}
