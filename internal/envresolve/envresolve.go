// Package envresolve resolves variable references within env maps,
// expanding values that reference other keys (e.g. BASE_URL=${SCHEME}://${HOST}).
package envresolve

import (
	"fmt"
	"regexp"
	"strings"
)

// MaxDepth limits recursive expansion to prevent infinite loops.
const MaxDepth = 10

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}`)

// Result holds the outcome of resolving a single key.
type Result struct {
	Key      string
	Original string
	Resolved string
	Expanded bool
	Err      error
}

// Resolve expands all variable references in the given env map.
// It returns a new map with resolved values and a slice of Results
// describing what changed or failed.
func Resolve(env map[string]string) (map[string]string, []Result) {
	out := make(map[string]string, len(env))
	results := make([]Result, 0, len(env))

	for k, v := range env {
		resolved, err := expand(v, env, 0)
		out[k] = resolved
		results = append(results, Result{
			Key:      k,
			Original: v,
			Resolved: resolved,
			Expanded: resolved != v,
			Err:      err,
		})
	}

	return out, results
}

// Violations returns only results that encountered an error during resolution.
func Violations(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.Err != nil {
			out = append(out, r)
		}
	}
	return out
}

func expand(value string, env map[string]string, depth int) (string, error) {
	if depth > MaxDepth {
		return value, fmt.Errorf("expansion depth exceeded for value %q", value)
	}

	var expandErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-1])
		replacement, ok := env[key]
		if !ok {
			expandErr = fmt.Errorf("undefined variable: %s", key)
			return match
		}
		expanded, err := expand(replacement, env, depth+1)
		if err != nil {
			expandErr = err
			return match
		}
		return expanded
	})

	return result, expandErr
}
