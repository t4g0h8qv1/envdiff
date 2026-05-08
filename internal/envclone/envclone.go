// Package envclone provides functionality to clone an env map into a new
// target file, optionally filtering keys by prefix or pattern.
package envclone

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

// Options controls how a clone operation is performed.
type Options struct {
	// KeyPrefix filters keys to only those starting with this prefix.
	KeyPrefix string
	// KeyContains filters keys to only those containing this substring.
	KeyContains string
	// Redact replaces values with empty strings in the output.
	Redact bool
	// OverwriteExisting allows overwriting the target file if it exists.
	OverwriteExisting bool
}

// Result summarises a completed clone operation.
type Result struct {
	Cloned  int
	Skipped int
	Target  string
}

// Clone writes a filtered, optionally redacted copy of src into the file at
// targetPath. It returns a Result describing what was written.
func Clone(src map[string]string, targetPath string, opts Options) (Result, error) {
	if !opts.OverwriteExisting {
		if _, err := os.Stat(targetPath); err == nil {
			return Result{}, fmt.Errorf("target file already exists: %s (use OverwriteExisting to force)", targetPath)
		}
	}

	keys := make([]string, 0, len(src))
	for k := range src {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	skipped := 0

	for _, k := range keys {
		if opts.KeyPrefix != "" && !strings.HasPrefix(k, opts.KeyPrefix) {
			skipped++
			continue
		}
		if opts.KeyContains != "" && !strings.Contains(k, opts.KeyContains) {
			skipped++
			continue
		}
		v := src[k]
		if opts.Redact {
			v = ""
		}
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}

	if err := os.WriteFile(targetPath, []byte(content), 0o644); err != nil {
		return Result{}, fmt.Errorf("writing clone to %s: %w", targetPath, err)
	}

	return Result{
		Cloned:  len(lines),
		Skipped: skipped,
		Target:  targetPath,
	}, nil
}
