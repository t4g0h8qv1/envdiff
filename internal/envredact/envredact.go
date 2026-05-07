// Package envredact provides utilities for redacting sensitive values
// in environment variable maps before display or export.
package envredact

import "strings"

// DefaultSensitivePatterns contains common key substrings that indicate
// a value should be redacted.
var DefaultSensitivePatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"APIKEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
	"ACCESS_KEY",
}

const redactedPlaceholder = "[REDACTED]"

// Redactor holds configuration for redaction.
type Redactor struct {
	Patterns []string
}

// New returns a Redactor using the provided patterns.
// If patterns is nil or empty, DefaultSensitivePatterns is used.
func New(patterns []string) *Redactor {
	if len(patterns) == 0 {
		patterns = DefaultSensitivePatterns
	}
	return &Redactor{Patterns: patterns}
}

// IsSensitive returns true if the given key matches any sensitive pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.Patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Redact returns a copy of the env map with sensitive values replaced
// by the redacted placeholder.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if r.IsSensitive(k) {
			out[k] = redactedPlaceholder
		} else {
			out[k] = v
		}
	}
	return out
}

// RedactAll applies redaction to multiple env maps and returns new copies.
func (r *Redactor) RedactAll(envs map[string]map[string]string) map[string]map[string]string {
	out := make(map[string]map[string]string, len(envs))
	for name, env := range envs {
		out[name] = r.Redact(env)
	}
	return out
}
