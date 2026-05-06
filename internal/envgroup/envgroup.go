// Package envgroup provides utilities for grouping and summarizing
// environment variables by prefix or custom category rules.
package envgroup

import (
	"sort"
	"strings"
)

// Group represents a named collection of env keys and their values.
type Group struct {
	Name string
	Keys map[string]string
}

// ByPrefix splits an env map into groups based on key prefixes.
// Keys are split on the first underscore; keys with no underscore
// are placed in the "_OTHER" group.
func ByPrefix(env map[string]string) []Group {
	buckets := make(map[string]map[string]string)

	for k, v := range env {
		prefix := "_OTHER"
		if idx := strings.Index(k, "_"); idx > 0 {
			prefix = k[:idx]
		}
		if buckets[prefix] == nil {
			buckets[prefix] = make(map[string]string)
		}
		buckets[prefix][k] = v
	}

	return sortedGroups(buckets)
}

// ByCategories groups env keys according to caller-supplied category
// definitions. Each category maps a label to a list of key prefixes.
// Keys not matched by any category land in "_OTHER".
func ByCategories(env map[string]string, categories map[string][]string) []Group {
	buckets := make(map[string]map[string]string)

	for k, v := range env {
		label := "_OTHER"
		for cat, prefixes := range categories {
			for _, p := range prefixes {
				if strings.HasPrefix(k, p) {
					label = cat
					break
				}
			}
			if label != "_OTHER" {
				break
			}
		}
		if buckets[label] == nil {
			buckets[label] = make(map[string]string)
		}
		buckets[label][k] = v
	}

	return sortedGroups(buckets)
}

// Summary returns a map of group name → number of keys in that group.
func Summary(groups []Group) map[string]int {
	out := make(map[string]int, len(groups))
	for _, g := range groups {
		out[g.Name] = len(g.Keys)
	}
	return out
}

func sortedGroups(buckets map[string]map[string]string) []Group {
	names := make([]string, 0, len(buckets))
	for n := range buckets {
		names = append(names, n)
	}
	sort.Strings(names)

	groups := make([]Group, 0, len(names))
	for _, n := range names {
		groups = append(groups, Group{Name: n, Keys: buckets[n]})
	}
	return groups
}
