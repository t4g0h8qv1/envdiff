package envdiff

import (
	"fmt"
	"strings"
)

// Severity represents the importance level of a drift item.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

// DriftItem describes a single detected drift between environments.
type DriftItem struct {
	Key      string
	EnvA     string
	EnvB     string
	ValueA   string
	ValueB   string
	Severity Severity
	Reason   string
}

// DriftReport holds all drift items found across a set of environments.
type DriftReport struct {
	Items []DriftItem
}

// HasDrift returns true if any drift items exist.
func (r *DriftReport) HasDrift() bool {
	return len(r.Items) > 0
}

// BySeverity returns items filtered to the given severity.
func (r *DriftReport) BySeverity(s Severity) []DriftItem {
	var out []DriftItem
	for _, item := range r.Items {
		if item.Severity == s {
			out = append(out, item)
		}
	}
	return out
}

// Detect compares two named environment maps and returns a DriftReport.
func Detect(envA, envB map[string]string, nameA, nameB string) *DriftReport {
	report := &DriftReport{}

	allKeys := mergeKeys(envA, envB)
	for _, key := range allKeys {
		valA, okA := envA[key]
		valB, okB := envB[key]

		switch {
		case okA && !okB:
			report.Items = append(report.Items, DriftItem{
				Key: key, EnvA: nameA, EnvB: nameB,
				ValueA: valA, Severity: SeverityWarning,
				Reason: fmt.Sprintf("key present in %s but missing in %s", nameA, nameB),
			})
		case !okA && okB:
			report.Items = append(report.Items, DriftItem{
				Key: key, EnvA: nameA, EnvB: nameB,
				ValueB: valB, Severity: SeverityWarning,
				Reason: fmt.Sprintf("key present in %s but missing in %s", nameB, nameA),
			})
		case valA != valB:
			sev := SeverityInfo
			if strings.HasPrefix(strings.ToUpper(key), "SECRET") || strings.HasPrefix(strings.ToUpper(key), "TOKEN") {
				sev = SeverityCritical
			}
			report.Items = append(report.Items, DriftItem{
				Key: key, EnvA: nameA, EnvB: nameB,
				ValueA: valA, ValueB: valB, Severity: sev,
				Reason: "values differ between environments",
			})
		}
	}
	return report
}

func mergeKeys(a, b map[string]string) []string {
	seen := make(map[string]struct{})
	var keys []string
	for k := range a {
		if _, ok := seen[k]; !ok {
			keys = append(keys, k)
			seen[k] = struct{}{}
		}
	}
	for k := range b {
		if _, ok := seen[k]; !ok {
			keys = append(keys, k)
			seen[k] = struct{}{}
		}
	}
	return keys
}
