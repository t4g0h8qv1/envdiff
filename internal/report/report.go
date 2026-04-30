package report

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the output format for the report.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Options holds configuration for rendering a report.
type Options struct {
	Format   Format
	Colorize bool
}

const (
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

// Render writes a human-readable diff report to w.
func Render(w io.Writer, results []diff.Result, leftName, rightName string, opts Options) error {
	if opts.Format == FormatJSON {
		return renderJSON(w, results, leftName, rightName)
	}
	return renderText(w, results, leftName, rightName, opts.Colorize)
}

func renderText(w io.Writer, results []diff.Result, leftName, rightName string, colorize bool) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "No differences found.")
		return err
	}

	fmt.Fprintf(w, "Comparing %s <-> %s\n", leftName, rightName)
	fmt.Fprintln(w, strings.Repeat("-", 40))

	for _, r := range results {
		switch r.Kind {
		case diff.Missing:
			prefix := fmt.Sprintf("[MISSING in %s]", r.MissingIn)
			if colorize {
				fmt.Fprintf(w, "%s%s%s KEY=%s\n", colorYellow, prefix, colorReset, r.Key)
			} else {
				fmt.Fprintf(w, "%s KEY=%s\n", prefix, r.Key)
			}
		case diff.Conflict:
			if colorize {
				fmt.Fprintf(w, "%s[CONFLICT]%s KEY=%s (%s=%q, %s=%q)\n",
					colorRed, colorReset, r.Key, leftName, r.LeftValue, rightName, r.RightValue)
			} else {
				fmt.Fprintf(w, "[CONFLICT] KEY=%s (%s=%q, %s=%q)\n",
					r.Key, leftName, r.LeftValue, rightName, r.RightValue)
			}
		}
	}
	return nil
}

func renderJSON(w io.Writer, results []diff.Result, leftName, rightName string) error {
	fmt.Fprintf(w, "{\n  \"left\": %q,\n  \"right\": %q,\n  \"differences\": [\n", leftName, rightName)
	for i, r := range results {
		comma := ","
		if i == len(results)-1 {
			comma = ""
		}
		fmt.Fprintf(w, "    {\"kind\": %q, \"key\": %q, \"left_value\": %q, \"right_value\": %q, \"missing_in\": %q}%s\n",
			r.Kind, r.Key, r.LeftValue, r.RightValue, r.MissingIn, comma)
	}
	fmt.Fprint(w, "  ]\n}\n")
	return nil
}
