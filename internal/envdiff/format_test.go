package envdiff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrint_TextNoDrift(t *testing.T) {
	r := &DriftReport{}
	var buf bytes.Buffer
	if err := Print(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestPrint_TextWithWarning(t *testing.T) {
	r := &DriftReport{
		Items: []DriftItem{
			{Key: "MISSING_KEY", EnvA: "prod", Severity: SeverityWarning,
				ValueA: "someval", Reason: "key present in prod but missing in staging"},
		},
	}
	var buf bytes.Buffer
	if err := Print(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARNING]") {
		t.Errorf("expected WARNING tag in output")
	}
	if !strings.Contains(out, "MISSING_KEY") {
		t.Errorf("expected key name in output")
	}
}

func TestPrint_TextWithCritical(t *testing.T) {
	r := &DriftReport{
		Items: []DriftItem{
			{Key: "SECRET_KEY", EnvA: "dev", EnvB: "prod",
				ValueA: "abc", ValueB: "xyz",
				Severity: SeverityCritical, Reason: "values differ between environments"},
		},
	}
	var buf bytes.Buffer
	if err := Print(&buf, r, FormatText); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[CRITICAL]") {
		t.Errorf("expected CRITICAL tag, got: %s", out)
	}
}

func TestPrint_JSONFormat(t *testing.T) {
	r := &DriftReport{
		Items: []DriftItem{
			{Key: "APP_ENV", EnvA: "dev", EnvB: "prod",
				ValueA: "development", ValueB: "production",
				Severity: SeverityInfo, Reason: "values differ between environments"},
		},
	}
	var buf bytes.Buffer
	if err := Print(&buf, r, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var items []DriftItem
	if err := json.Unmarshal(buf.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if len(items) != 1 || items[0].Key != "APP_ENV" {
		t.Errorf("unexpected JSON content: %+v", items)
	}
}

func TestPrint_JSONNoDrift(t *testing.T) {
	r := &DriftReport{}
	var buf bytes.Buffer
	if err := Print(&buf, r, FormatJSON); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "null") && !strings.Contains(buf.String(), "[]") {
		t.Errorf("expected empty JSON array or null, got: %s", buf.String())
	}
}
