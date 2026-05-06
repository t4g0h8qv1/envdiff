// Package audit provides change tracking for .env file comparisons.
// It records diff results with timestamps and metadata for audit trails.
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time     `json:"timestamp"`
	Files     []string      `json:"files"`
	Results   []diff.Result `json:"results"`
	Summary   Summary       `json:"summary"`
}

// Summary holds aggregated counts from a diff run.
type Summary struct {
	Total     int `json:"total"`
	Missing   int `json:"missing"`
	Conflicts int `json:"conflicts"`
}

// NewEntry creates a new audit Entry from diff results and file paths.
func NewEntry(files []string, results []diff.Result) Entry {
	s := Summary{Total: len(results)}
	for _, r := range results {
		switch r.Status {
		case "missing":
			s.Missing++
		case "conflict":
			s.Conflicts++
		}
	}
	return Entry{
		Timestamp: time.Now().UTC(),
		Files:     files,
		Results:   results,
		Summary:   s,
	}
}

// Append writes an audit entry to the given log file in JSON Lines format.
func Append(path string, entry Entry) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll reads all audit entries from a JSON Lines log file.
func ReadAll(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("audit: read log file: %w", err)
	}

	var entries []Entry
	dec := json.NewDecoder(bytes.NewReader(data))
	for dec.More() {
		var e Entry
		if err := dec.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}
