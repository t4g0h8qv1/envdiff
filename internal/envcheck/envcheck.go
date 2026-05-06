// Package envcheck verifies that all keys defined in a reference .env file
// are present (and optionally non-empty) in one or more target env maps.
package envcheck

import "fmt"

// Result holds the outcome of a single key check across a named environment.
type Result struct {
	Key     string
	EnvName string
	Missing bool
	Empty   bool
}

// Options controls how the check behaves.
type Options struct {
	// RequireNonEmpty treats keys that exist but have an empty value as failures.
	RequireNonEmpty bool
}

// Check verifies that every key present in reference exists in each of the
// provided target maps. The map keys of targets are used as environment names
// in the returned results.
func Check(reference map[string]string, targets map[string]map[string]string, opts Options) []Result {
	var results []Result

	for key := range reference {
		for envName, envMap := range targets {
			val, exists := envMap[key]
			if !exists {
				results = append(results, Result{
					Key:     key,
					EnvName: envName,
					Missing: true,
				})
				continue
			}
			if opts.RequireNonEmpty && val == "" {
				results = append(results, Result{
					Key:     key,
					EnvName: envName,
					Empty:   true,
				})
			}
		}
	}

	return results
}

// Format returns a human-readable summary of the check results.
func Format(results []Result) string {
	if len(results) == 0 {
		return "All keys present and valid."
	}
	out := ""
	for _, r := range results {
		switch {
		case r.Missing:
			out += fmt.Sprintf("[MISSING] key %q not found in %q\n", r.Key, r.EnvName)
		case r.Empty:
			out += fmt.Sprintf("[EMPTY]   key %q is empty in %q\n", r.Key, r.EnvName)
		}
	}
	return out
}
