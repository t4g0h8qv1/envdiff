// Package envpin provides functionality to "pin" specific env keys to
// expected values and detect when those values drift from the pinned state.
package envpin

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"
)

// Pin represents a single pinned key-value expectation.
type Pin struct {
	Key      string `json:"key"`
	Expected string `json:"expected"`
	PinnedAt string `json:"pinned_at"`
}

// PinFile is the top-level structure stored on disk.
type PinFile struct {
	Pins []Pin `json:"pins"`
}

// Violation describes a key whose current value differs from its pinned value.
type Violation struct {
	Key      string
	Pinned   string
	Actual   string
	Missing  bool
}

// Save writes a set of key pins derived from env to the given path.
func Save(path string, env map[string]string, keys []string) error {
	pf := PinFile{}
	for _, k := range keys {
		v, ok := env[k]
		if !ok {
			return fmt.Errorf("key %q not found in env map", k)
		}
		pf.Pins = append(pf.Pins, Pin{
			Key:      k,
			Expected: v,
			PinnedAt: time.Now().UTC().Format(time.RFC3339),
		})
	}
	sort.Slice(pf.Pins, func(i, j int) bool { return pf.Pins[i].Key < pf.Pins[j].Key })
	data, err := json.MarshalIndent(pf, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal pins: %w", err)
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a PinFile from disk.
func Load(path string) (PinFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PinFile{}, fmt.Errorf("read pin file: %w", err)
	}
	var pf PinFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return PinFile{}, fmt.Errorf("parse pin file: %w", err)
	}
	return pf, nil
}

// Check compares the current env map against a PinFile and returns any violations.
func Check(pf PinFile, env map[string]string) []Violation {
	var violations []Violation
	for _, pin := range pf.Pins {
		actual, ok := env[pin.Key]
		if !ok {
			violations = append(violations, Violation{Key: pin.Key, Pinned: pin.Expected, Missing: true})
			continue
		}
		if actual != pin.Expected {
			violations = append(violations, Violation{Key: pin.Key, Pinned: pin.Expected, Actual: actual})
		}
	}
	return violations
}
