package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a linting rule applied to env file keys and values.
type Rule string

const (
	RuleUppercaseKeys   Rule = "uppercase_keys"
	RuleNoSpacesInKeys  Rule = "no_spaces_in_keys"
	RuleNoEmptyValues   Rule = "no_empty_values"
	RuleKeyFormat       Rule = "key_format"
)

// Violation represents a single lint issue found in an env map.
type Violation struct {
	File    string
	Key     string
	Rule    Rule
	Message string
}

var validKeyRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// Check runs all enabled lint rules against the provided env maps.
// envMaps is a map of filename -> key/value pairs.
func Check(envMaps map[string]map[string]string, rules []Rule) []Violation {
	enabledRules := make(map[Rule]bool, len(rules))
	for _, r := range rules {
		enabledRules[r] = true
	}

	var violations []Violation

	for file, env := range envMaps {
		for key, value := range env {
			if enabledRules[RuleUppercaseKeys] {
				if key != strings.ToUpper(key) {
					violations = append(violations, Violation{
						File:    file,
						Key:     key,
						Rule:    RuleUppercaseKeys,
						Message: fmt.Sprintf("key %q should be uppercase", key),
					})
				}
			}
			if enabledRules[RuleNoSpacesInKeys] {
				if strings.Contains(key, " ") {
					violations = append(violations, Violation{
						File:    file,
						Key:     key,
						Rule:    RuleNoSpacesInKeys,
						Message: fmt.Sprintf("key %q contains spaces", key),
					})
				}
			}
			if enabledRules[RuleNoEmptyValues] {
				if strings.TrimSpace(value) == "" {
					violations = append(violations, Violation{
						File:    file,
						Key:     key,
						Rule:    RuleNoEmptyValues,
						Message: fmt.Sprintf("key %q has an empty value", key),
					})
				}
			}
			if enabledRules[RuleKeyFormat] {
				if !validKeyRe.MatchString(key) {
					violations = append(violations, Violation{
						File:    file,
						Key:     key,
						Rule:    RuleKeyFormat,
						Message: fmt.Sprintf("key %q does not match pattern [A-Z][A-Z0-9_]*", key),
					})
				}
			}
		}
	}

	return violations
}
