package report_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/envdiff/internal/diff"
	"github.com/user/envdiff/internal/report"
)

func TestRender_NoDifferences(t *testing.T) {
	var buf bytes.Buffer
	err := report.Render(&buf, nil, "dev.env", "prod.env", report.Options{Format: report.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No differences found") {
		t.Errorf("expected no-differences message, got: %s", buf.String())
	}
}

func TestRender_MissingKey(t *testing.T) {
	results := []diff.Result{
		{Kind: diff.Missing, Key: "SECRET_KEY", MissingIn: "prod.env"},
	}
	var buf bytes.Buffer
	err := report.Render(&buf, results, "dev.env", "prod.env", report.Options{Format: report.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "MISSING") {
		t.Errorf("expected MISSING label, got: %s", out)
	}
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected key name in output, got: %s", out)
	}
}

func TestRender_Conflict(t *testing.T) {
	results := []diff.Result{
		{Kind: diff.Conflict, Key: "DB_HOST", LeftValue: "localhost", RightValue: "db.prod.internal"},
	}
	var buf bytes.Buffer
	err := report.Render(&buf, results, "dev.env", "prod.env", report.Options{Format: report.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "CONFLICT") {
		t.Errorf("expected CONFLICT label, got: %s", out)
	}
	if !strings.Contains(out, "localhost") || !strings.Contains(out, "db.prod.internal") {
		t.Errorf("expected both values in output, got: %s", out)
	}
}

func TestRender_JSONFormat(t *testing.T) {
	results := []diff.Result{
		{Kind: diff.Missing, Key: "API_KEY", MissingIn: "prod.env"},
	}
	var buf bytes.Buffer
	err := report.Render(&buf, results, "dev.env", "prod.env", report.Options{Format: report.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "\"key\"") {
		t.Errorf("expected JSON key field, got: %s", out)
	}
	if !strings.Contains(out, "API_KEY") {
		t.Errorf("expected API_KEY in JSON output, got: %s", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(out), "{") {
		t.Errorf("expected JSON object, got: %s", out)
	}
}
