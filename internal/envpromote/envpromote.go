// Package envpromote provides utilities for promoting env vars from one
// environment to another (e.g. staging → production), applying promotion
// rules such as key allowlists, value transforms, and conflict policies.
package envpromote

import (
	"fmt"
	"strings"
)

// Policy controls how conflicts are handled during promotion.
type Policy string

const (
	// PolicySkip leaves the destination value unchanged on conflict.
	PolicySkip Policy = "skip"
	// PolicyOverwrite replaces the destination value with the source value.
	PolicyOverwrite Policy = "overwrite"
	// PolicyError returns an error when a conflict is detected.
	PolicyError Policy = "error"
)

// Options configures a promotion run.
type Options struct {
	// AllowKeys restricts promotion to only these keys. Empty means all keys.
	AllowKeys []string
	// DenyKeys excludes these keys from promotion.
	DenyKeys []string
	// Policy determines conflict resolution behaviour.
	Policy Policy
}

// Result captures what happened to a single key during promotion.
type Result struct {
	Key      string
	Action   string // "promoted", "skipped", "overwritten"
	OldValue string
	NewValue string
}

// Promote copies keys from src into dst according to opts.
// It returns a slice of Result describing every key that was considered.
func Promote(src, dst map[string]string, opts Options) ([]Result, error) {
	allowSet := toSet(opts.AllowKeys)
	denySet := toSet(opts.DenyKeys)

	var results []Result

	for k, srcVal := range src {
		if len(allowSet) > 0 && !allowSet[k] {
			continue
		}
		if denySet[k] {
			continue
		}

		dstVal, exists := dst[k]

		switch {
		case !exists:
			dst[k] = srcVal
			results = append(results, Result{Key: k, Action: "promoted", NewValue: srcVal})

		case srcVal == dstVal:
			results = append(results, Result{Key: k, Action: "skipped", OldValue: dstVal, NewValue: srcVal})

		default:
			switch opts.Policy {
			case PolicyOverwrite:
				dst[k] = srcVal
				results = append(results, Result{Key: k, Action: "overwritten", OldValue: dstVal, NewValue: srcVal})
			case PolicyError:
				return nil, fmt.Errorf("conflict on key %q: src=%q dst=%q", k, srcVal, dstVal)
			default: // PolicySkip
				results = append(results, Result{Key: k, Action: "skipped", OldValue: dstVal, NewValue: srcVal})
			}
		}
	}

	return results, nil
}

// Summary returns a human-readable summary line for a slice of results.
func Summary(results []Result) string {
	counts := map[string]int{}
	for _, r := range results {
		counts[r.Action]++
	}
	parts := []string{}
	for _, action := range []string{"promoted", "overwritten", "skipped"} {
		if n := counts[action]; n > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", n, action))
		}
	}
	if len(parts) == 0 {
		return "nothing to promote"
	}
	return strings.Join(parts, ", ")
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
