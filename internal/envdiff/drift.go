package envdiff

import (
	"fmt"
	"sort"
	"strings"
)

// DriftReport summarizes drift across multiple named environments.
type DriftReport struct {
	Environments []string
	Entries      []DriftEntry
}

// DriftEntry represents a key that differs across environments.
type DriftEntry struct {
	Key    string
	Values map[string]string // env name -> value (empty string means missing)
	Status string            // "missing", "conflict", "ok"
}

// DetectMulti compares more than two environments and returns a DriftReport.
func DetectMulti(envs map[string]map[string]string) DriftReport {
	allKeys := mergeAllKeys(envs)
	envNames := sortedEnvNames(envs)

	var entries []DriftEntry
	for _, key := range allKeys {
		values := make(map[string]string, len(envNames))
		for _, name := range envNames {
			v, ok := envs[name][key]
			if ok {
				values[name] = v
			} else {
				values[name] = ""
			}
		}
		status := classifyDrift(values, envs, key)
		if status != "ok" {
			entries = append(entries, DriftEntry{Key: key, Values: values, Status: status})
		}
	}

	return DriftReport{Environments: envNames, Entries: entries}
}

func classifyDrift(values map[string]string, envs map[string]map[string]string, key string) string {
	present := 0
	for _, env := range envs {
		if _, ok := env[key]; ok {
			present++
		}
	}
	if present < len(envs) {
		return "missing"
	}
	first := ""
	for _, v := range values {
		if first == "" {
			first = v
		} else if v != first {
			return "conflict"
		}
	}
	return "ok"
}

func mergeAllKeys(envs map[string]map[string]string) []string {
	seen := map[string]struct{}{}
	for _, m := range envs {
		for k := range m {
			seen[k] = struct{}{}
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func sortedEnvNames(envs map[string]map[string]string) []string {
	names := make([]string, 0, len(envs))
	for n := range envs {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Summary returns a human-readable one-line summary of the drift report.
func (r DriftReport) Summary() string {
	if len(r.Entries) == 0 {
		return fmt.Sprintf("no drift detected across %s", strings.Join(r.Environments, ", "))
	}
	missing, conflicts := 0, 0
	for _, e := range r.Entries {
		switch e.Status {
		case "missing":
			missing++
		case "conflict":
			conflicts++
		}
	}
	return fmt.Sprintf("%d missing, %d conflicting keys across %s",
		missing, conflicts, strings.Join(r.Environments, ", "))
}
