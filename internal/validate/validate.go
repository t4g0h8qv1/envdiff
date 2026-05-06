package validate

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a validation rule for an environment variable key or value.
type Rule struct {
	Key     string
	Pattern string
	Required bool
}

// Violation describes a validation failure.
type Violation struct {
	File    string
	Key     string
	Message string
}

// Check validates the parsed env maps against a set of rules.
// It returns a slice of violations found across all provided files.
func Check(envMaps map[string]map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, rule := range rules {
		for file, env := range envMaps {
			val, exists := env[rule.Key]

			if rule.Required && !exists {
				violations = append(violations, Violation{
					File:    file,
					Key:     rule.Key,
					Message: "required key is missing",
				})
				continue
			}

			if !exists {
				continue
			}

			if rule.Pattern != "" {
				matched, err := regexp.MatchString(rule.Pattern, val)
				if err != nil {
					violations = append(violations, Violation{
						File:    file,
						Key:     rule.Key,
						Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
					})
					continue
				}
				if !matched {
					violations = append(violations, Violation{
						File:    file,
						Key:     rule.Key,
						Message: fmt.Sprintf("value %q does not match pattern %q", val, rule.Pattern),
					})
				}
			}

			if strings.TrimSpace(val) == "" && rule.Required {
				violations = append(violations, Violation{
					File:    file,
					Key:     rule.Key,
					Message: "required key has empty value",
				})
			}
		}
	}

	return violations
}
