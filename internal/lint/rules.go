package lint

// AllRules returns the full set of available lint rules.
func AllRules() []Rule {
	return []Rule{
		RuleUppercaseKeys,
		RuleNoSpacesInKeys,
		RuleNoEmptyValues,
		RuleKeyFormat,
	}
}

// DefaultRules returns the recommended set of lint rules for most projects.
func DefaultRules() []Rule {
	return []Rule{
		RuleUppercaseKeys,
		RuleNoSpacesInKeys,
		RuleKeyFormat,
	}
}

// ParseRule converts a string to a Rule, returning an error if unknown.
func ParseRule(s string) (Rule, bool) {
	switch Rule(s) {
	case RuleUppercaseKeys,
		RuleNoSpacesInKeys,
		RuleNoEmptyValues,
		RuleKeyFormat:
		return Rule(s), true
	}
	return "", false
}

// RuleDescriptions returns a human-readable description for each rule.
func RuleDescriptions() map[Rule]string {
	return map[Rule]string{
		RuleUppercaseKeys:  "All keys must be fully uppercase",
		RuleNoSpacesInKeys: "Keys must not contain spaces",
		RuleNoEmptyValues:  "Values must not be empty or whitespace-only",
		RuleKeyFormat:      "Keys must match pattern [A-Z][A-Z0-9_]*",
	}
}
