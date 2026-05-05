package filter

import (
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Options holds filtering criteria for diff results.
type Options struct {
	OnlyMissing  bool
	OnlyConflicts bool
	KeyPrefix    string
	KeyContains  string
}

// Apply filters a slice of diff.Result entries based on the provided Options.
func Apply(results []diff.Result, opts Options) []diff.Result {
	var filtered []diff.Result

	for _, r := range results {
		if opts.OnlyMissing && r.Type != diff.Missing {
			continue
		}
		if opts.OnlyConflicts && r.Type != diff.Conflict {
			continue
		}
		if opts.KeyPrefix != "" && !strings.HasPrefix(r.Key, opts.KeyPrefix) {
			continue
		}
		if opts.KeyContains != "" && !strings.Contains(r.Key, opts.KeyContains) {
			continue
		}
		filtered = append(filtered, r)
	}

	return filtered
}
