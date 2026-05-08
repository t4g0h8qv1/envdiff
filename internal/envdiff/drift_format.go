package envdiff

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
)

// PrintDrift writes a formatted drift report to w.
func PrintDrift(w io.Writer, report DriftReport, format string) error {
	switch strings.ToLower(format) {
	case "json":
		return printDriftJSON(w, report)
	default:
		return printDriftText(w, report)
	}
}

func printDriftText(w io.Writer, report DriftReport) error {
	fmt.Fprintf(w, "Environments: %s\n", strings.Join(report.Environments, " | "))
	fmt.Fprintf(w, "%-40s %-12s", "KEY", "STATUS")
	for _, env := range report.Environments {
		fmt.Fprintf(w, " %-20s", env)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, strings.Repeat("-", 40+12+len(report.Environments)*21))

	if len(report.Entries) == 0 {
		fmt.Fprintln(w, "  No drift detected.")
		return nil
	}

	for _, e := range report.Entries {
		color := colorYellow
		if e.Status == "conflict" {
			color = colorRed
		}
		fmt.Fprintf(w, "%s%-40s %-12s%s", color, e.Key, e.Status, colorReset)
		for _, env := range report.Environments {
			v := e.Values[env]
			if v == "" {
				v = "(missing)"
			}
			if len(v) > 18 {
				v = v[:15] + "..."
			}
			fmt.Fprintf(w, " %-20s", v)
		}
		fmt.Fprintln(w)
	}
	return nil
}

func printDriftJSON(w io.Writer, report DriftReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
