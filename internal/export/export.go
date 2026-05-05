package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/user/envdiff/internal/diff"
)

// Format represents the export format type.
type Format string

const (
	FormatCSV      Format = "csv"
	FormatMarkdown Format = "markdown"
	FormatJSON     Format = "json"
)

// Write exports the diff results to the given writer in the specified format.
func Write(w io.Writer, results []diff.Result, format Format) error {
	switch format {
	case FormatCSV:
		return writeCSV(w, results)
	case FormatMarkdown:
		return writeMarkdown(w, results)
	case FormatJSON:
		return writeJSON(w, results)
	default:
		return fmt.Errorf("unsupported export format: %s", format)
	}
}

func writeCSV(w io.Writer, results []diff.Result) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"key", "status", "left_value", "right_value"}); err != nil {
		return err
	}
	for _, r := range results {
		row := []string{r.Key, string(r.Status), r.LeftValue, r.RightValue}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func writeMarkdown(w io.Writer, results []diff.Result) error {
	if _, err := fmt.Fprintln(w, "| Key | Status | Left Value | Right Value |"); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(w, "|-----|--------|------------|-------------|"); err != nil {
		return err
	}
	for _, r := range results {
		line := fmt.Sprintf("| %s | %s | %s | %s |",
			r.Key, string(r.Status),
			escapeMD(r.LeftValue), escapeMD(r.RightValue))
		if _, err := fmt.Fprintln(w, line); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, results []diff.Result) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func escapeMD(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
