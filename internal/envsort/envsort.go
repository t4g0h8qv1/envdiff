// Package envsort provides utilities for sorting and normalizing env maps.
package envsort

import (
	"sort"
	"strings"
)

// SortedKeys returns the keys of the given map in sorted order.
func SortedKeys(env map[string]string) []string {
	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Normalize returns a new map with all keys converted to uppercase.
func Normalize(env map[string]string) map[string]string {
	normalized := make(map[string]string, len(env))
	for k, v := range env {
		normalized[strings.ToUpper(k)] = v
	}
	return normalized
}

// SortedPairs returns key=value pairs from the map in sorted key order.
func SortedPairs(env map[string]string) []string {
	keys := SortedKeys(env)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+env[k])
	}
	return pairs
}

// GroupByPrefix groups keys by their prefix (up to the first underscore).
// Keys without an underscore are grouped under "".
func GroupByPrefix(env map[string]string) map[string][]string {
	groups := make(map[string][]string)
	for k := range env {
		parts := strings.SplitN(k, "_", 2)
		prefix := ""
		if len(parts) == 2 {
			prefix = parts[0]
		}
		groups[prefix] = append(groups[prefix], k)
	}
	for prefix := range groups {
		sort.Strings(groups[prefix])
	}
	return groups
}
