// Package suggest provides recommendations for resolving diff results.
package suggest

import (
	"fmt"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Suggestion represents a recommended action for a diff result.
type Suggestion struct {
	Key    string
	Kind   string // "missing", "conflict"
	File   string
	Action string
}

// Generate returns a list of suggestions based on diff results.
func Generate(results []diff.Result) []Suggestion {
	var suggestions []Suggestion

	for _, r := range results {
		switch r.Kind {
		case "missing":
			suggestions = append(suggestions, Suggestion{
				Key:    r.Key,
				Kind:   "missing",
				File:   r.File,
				Action: fmt.Sprintf("Add key %q to %s", r.Key, r.File),
			})
		case "conflict":
			files := strings.Join(r.Files, ", ")
			suggestions = append(suggestions, Suggestion{
				Key:    r.Key,
				Kind:   "conflict",
				File:   files,
				Action: fmt.Sprintf("Resolve conflicting values for key %q across: %s", r.Key, files),
			})
		}
	}

	return suggestions
}

// Format returns a human-readable string of all suggestions.
func Format(suggestions []Suggestion) string {
	if len(suggestions) == 0 {
		return "No suggestions — all environments look consistent.\n"
	}

	var sb strings.Builder
	sb.WriteString("Suggestions:\n")
	for i, s := range suggestions {
		sb.WriteString(fmt.Sprintf("  %d. [%s] %s\n", i+1, strings.ToUpper(s.Kind), s.Action))
	}
	return sb.String()
}
