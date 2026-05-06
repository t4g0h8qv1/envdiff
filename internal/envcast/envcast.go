// Package envcast provides utilities for casting env map values to typed Go values.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds the typed value result for a single key.
type Result struct {
	Key   string
	Raw   string
	Error error
}

// CastMap attempts to cast values in the env map according to the provided type hints.
// typeHints maps key names to desired types: "int", "bool", "float", "string".
func CastMap(env map[string]string, typeHints map[string]string) []Result {
	results := make([]Result, 0, len(typeHints))
	for key, kind := range typeHints {
		raw, ok := env[key]
		if !ok {
			results = append(results, Result{Key: key, Raw: "", Error: fmt.Errorf("key %q not found", key)})
			continue
		}
		var castErr error
		switch strings.ToLower(kind) {
		case "int":
			_, castErr = strconv.Atoi(raw)
		case "bool":
			_, castErr = strconv.ParseBool(raw)
		case "float":
			_, castErr = strconv.ParseFloat(raw, 64)
		case "string":
			// always valid
		default:
			castErr = fmt.Errorf("unknown type hint %q for key %q", kind, key)
		}
		if castErr != nil {
			castErr = fmt.Errorf("key %q: cannot cast %q to %s: %w", key, raw, kind, castErr)
		}
		results = append(results, Result{Key: key, Raw: raw, Error: castErr})
	}
	return results
}

// Violations returns only the Results that have errors.
func Violations(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Error != nil {
			out = append(out, r)
		}
	}
	return out
}
