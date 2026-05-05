package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/envdiff/internal/parser"
)

// EnvFile represents a loaded environment file with its name and parsed entries.
type EnvFile struct {
	Name    string
	Path    string
	Entries map[string]string
}

// LoadFiles loads one or more .env files from the given paths.
// It returns a slice of EnvFile or an error if any file cannot be read or parsed.
func LoadFiles(paths []string) ([]EnvFile, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("no file paths provided")
	}

	var files []EnvFile
	for _, p := range paths {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %s", p)
		}

		entries, err := parser.ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %w", p, err)
		}

		files = append(files, EnvFile{
			Name:    filepath.Base(p),
			Path:    p,
			Entries: entries,
		})
	}

	return files, nil
}

// LoadDir scans a directory for files matching the pattern *.env or .env*
// and loads all of them.
func LoadDir(dir string) ([]EnvFile, error) {
	patterns := []string{
		filepath.Join(dir, "*.env"),
		filepath.Join(dir, ".env*"),
	}

	var paths []string
	seen := map[string]bool{}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("glob error for pattern %s: %w", pattern, err)
		}
		for _, m := range matches {
			if !seen[m] {
				seen[m] = true
				paths = append(paths, m)
			}
		}
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("no .env files found in directory: %s", dir)
	}

	return LoadFiles(paths)
}
