// Package envnorm normalizes .env file keys and values according to
// configurable conventions (e.g. uppercase keys, trimmed values, no empty keys).
package envnorm

import (
	"fmt"
	"strings"
)

// Options controls which normalization rules are applied.
type Options struct {
	UppercaseKeys  bool
	TrimValues     bool
	RemoveEmpty    bool
	QuoteValues    bool
}

// DefaultOptions returns a sensible default normalization configuration.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys: true,
		TrimValues:    true,
		RemoveEmpty:   false,
		QuoteValues:   false,
	}
}

// Violation describes a key/value pair that was changed during normalization.
type Violation struct {
	Key      string
	Original string
	Normalized string
	Reason   string
}

// Normalize applies the given options to the provided env map and returns a
// new normalized map along with a list of changes made.
func Normalize(env map[string]string, opts Options) (map[string]string, []Violation) {
	result := make(map[string]string, len(env))
	var violations []Violation

	for k, v := range env {
		newKey := k
		newVal := v

		if opts.UppercaseKeys {
			upper := strings.ToUpper(k)
			if upper != k {
				violations = append(violations, Violation{
					Key:        k,
					Original:   k,
					Normalized: upper,
					Reason:     "key uppercased",
				})
				newKey = upper
			}
		}

		if opts.TrimValues {
			trimmed := strings.TrimSpace(v)
			if trimmed != v {
				violations = append(violations, Violation{
					Key:        newKey,
					Original:   v,
					Normalized: trimmed,
					Reason:     "value trimmed",
				})
				newVal = trimmed
			}
		}

		if opts.RemoveEmpty && newVal == "" {
			violations = append(violations, Violation{
				Key:        newKey,
				Original:   newVal,
				Normalized: "",
				Reason:     "empty value removed",
			})
			continue
		}

		if opts.QuoteValues && newVal != "" && !isQuoted(newVal) {
			quoted := fmt.Sprintf("%q", newVal)
			violations = append(violations, Violation{
				Key:        newKey,
				Original:   newVal,
				Normalized: quoted,
				Reason:     "value quoted",
			})
			newVal = quoted
		}

		result[newKey] = newVal
	}

	return result, violations
}

func isQuoted(s string) bool {
	return (strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) ||
		(strings.HasPrefix(s, "'") && strings.HasSuffix(s, "'"))
}
