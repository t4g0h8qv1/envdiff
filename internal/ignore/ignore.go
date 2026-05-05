package ignore

import (
	"bufio"
	"os"
	"strings"
)

// Rules holds a set of key patterns to ignore during comparison.
type Rules struct {
	keys map[string]struct{}
	prefixes []string
}

// LoadFile reads an ignore file where each line is a key or prefix pattern.
// Lines starting with '#' are treated as comments. Prefix patterns end with '*'.
func LoadFile(path string) (*Rules, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rules := &Rules{
		keys: make(map[string]struct{}),
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasSuffix(line, "*") {
			rules.prefixes = append(rules.prefixes, strings.TrimSuffix(line, "*"))
		} else {
			rules.keys[line] = struct{}{}
		}
	}

	return rules, scanner.Err()
}

// NewRules creates an empty Rules set.
func NewRules() *Rules {
	return &Rules{keys: make(map[string]struct{})}
}

// ShouldIgnore returns true if the given key matches any ignore rule.
func (r *Rules) ShouldIgnore(key string) bool {
	if _, ok := r.keys[key]; ok {
		return true
	}
	for _, prefix := range r.prefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// FilterKeys removes ignored keys from the provided map, returning a new map.
func (r *Rules) FilterKeys(env map[string]string) map[string]string {
	filtered := make(map[string]string, len(env))
	for k, v := range env {
		if !r.ShouldIgnore(k) {
			filtered[k] = v
		}
	}
	return filtered
}
