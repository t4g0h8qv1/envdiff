package envdiff

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format controls the output format for a DriftReport.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Print writes the DriftReport to w in the requested format.
func Print(w io.Writer, r *DriftReport, format Format) error {
	switch format {
	case FormatJSON:
		return printJSON(w, r)
	default:
		return printText(w, r)
	}
}

func printText(w io.Writer, r *DriftReport) error {
	if !r.HasDrift() {
		_, err := fmt.Fprintln(w, "No drift detected.")
		return err
	}
	for _, item := range r.Items {
		severityTag := strings.ToUpper(string(item.Severity))
		switch item.Severity {
		case SeverityCritical:
			_, err := fmt.Fprintf(w, "[%s] %s: %s\n  %s=%q  |  %s=%q\n",
				severityTag, item.Key, item.Reason, item.EnvA, item.ValueA, item.EnvB, item.ValueB)
			if err != nil {
				return err
			}
		case SeverityWarning:
			env, val := item.EnvA, item.ValueA
			if val == "" {
				env, val = item.EnvB, item.ValueB
			}
			_, err := fmt.Fprintf(w, "[%s] %s: %s\n  present in %s=%q\n",
				severityTag, item.Key, item.Reason, env, val)
			if err != nil {
				return err
			}
		default:
			_, err := fmt.Fprintf(w, "[%s] %s: %s\n  %s=%q  |  %s=%q\n",
				severityTag, item.Key, item.Reason, item.EnvA, item.ValueA, item.EnvB, item.ValueB)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printJSON(w io.Writer, r *DriftReport) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r.Items)
}
