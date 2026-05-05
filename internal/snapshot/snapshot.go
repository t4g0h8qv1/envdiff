package snapshot

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/user/envdiff/internal/diff"
)

// Snapshot represents a saved diff result at a point in time.
type Snapshot struct {
	CreatedAt time.Time    `json:"created_at"`
	Label     string       `json:"label"`
	Results   []diff.Result `json:"results"`
}

// Save writes a snapshot of the given diff results to a JSON file.
func Save(path, label string, results []diff.Result) error {
	s := Snapshot{
		CreatedAt: time.Now().UTC(),
		Label:     label,
		Results:   results,
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal failed: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("snapshot: write failed: %w", err)
	}

	return nil
}

// Load reads a snapshot from a JSON file.
func Load(path string) (*Snapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot: read failed: %w", err)
	}

	var s Snapshot
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("snapshot: unmarshal failed: %w", err)
	}

	return &s, nil
}

// Compare returns results that are new or changed compared to a baseline snapshot.
func Compare(baseline *Snapshot, current []diff.Result) []diff.Result {
	baselineKeys := make(map[string]diff.Result, len(baseline.Results))
	for _, r := range baseline.Results {
		baselineKeys[r.Key] = r
	}

	var changed []diff.Result
	for _, r := range current {
		if prev, ok := baselineKeys[r.Key]; !ok || prev.Type != r.Type {
			changed = append(changed, r)
		}
	}

	return changed
}
