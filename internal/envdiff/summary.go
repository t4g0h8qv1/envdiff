package envdiff

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// SummaryOptions controls what is included in the summary report.
type SummaryOptions struct {
	ShowCounts   bool
	ShowSeverity bool
	ShowKeys     bool
}

// DefaultSummaryOptions returns sensible defaults.
func DefaultSummaryOptions() SummaryOptions {
	return SummaryOptions{
		ShowCounts:   true,
		ShowSeverity: true,
		ShowKeys:     false,
	}
}

// Summary produces a human-readable summary of a DriftReport written to w.
func Summary(w io.Writer, report DriftReport, opts SummaryOptions) {
	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "✔ No drift detected across all environments.")
		return
	}

	counts := map[string]int{}
	bySeverity := map[string][]string{}

	for _, e := range report.Entries {
		sev := string(e.Severity)
		counts[sev]++
		bySeverity[sev] = append(bySeverity[sev], e.Key)
	}

	total := len(report.Entries)
	fmt.Fprintf(w, "Drift summary: %d issue(s) found across %d environment(s).\n", total, len(report.Environments))

	if opts.ShowCounts || opts.ShowSeverity {
		severities := []string{"critical", "warning", "info"}
		for _, sev := range severities {
			count, ok := counts[sev]
			if !ok {
				continue
			}
			line := fmt.Sprintf("  %-10s %d", strings.ToUpper(sev)+":", count)
			if opts.ShowKeys {
				keys := bySeverity[sev]
				sort.Strings(keys)
				line += fmt.Sprintf(" (%s)", strings.Join(keys, ", "))
			}
			fmt.Fprintln(w, line)
		}
	}
}
