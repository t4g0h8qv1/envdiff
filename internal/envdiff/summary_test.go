package envdiff

import (
	"strings"
	"testing"
)

func sampleReport() DriftReport {
	return DriftReport{
		Environments: []string{"dev", "prod", "staging"},
		Entries: []DriftEntry{
			{Key: "SECRET_KEY", Severity: SeverityCritical, Kind: "conflict"},
			{Key: "DB_HOST", Severity: SeverityWarning, Kind: "missing"},
			{Key: "APP_ENV", Severity: SeverityInfo, Kind: "conflict"},
			{Key: "LOG_LEVEL", Severity: SeverityWarning, Kind: "missing"},
		},
	}
}

func TestSummary_NoDrift(t *testing.T) {
	var b strings.Builder
	Summary(&b, DriftReport{}, DefaultSummaryOptions())
	out := b.String()
	if !strings.Contains(out, "No drift detected") {
		t.Errorf("expected no-drift message, got: %s", out)
	}
}

func TestSummary_ShowsCounts(t *testing.T) {
	var b strings.Builder
	Summary(&b, sampleReport(), DefaultSummaryOptions())
	out := b.String()
	if !strings.Contains(out, "4 issue(s)") {
		t.Errorf("expected total count, got: %s", out)
	}
	if !strings.Contains(out, "CRITICAL:") {
		t.Errorf("expected CRITICAL line, got: %s", out)
	}
	if !strings.Contains(out, "WARNING:") {
		t.Errorf("expected WARNING line, got: %s", out)
	}
	if !strings.Contains(out, "INFO:") {
		t.Errorf("expected INFO line, got: %s", out)
	}
}

func TestSummary_ShowKeys(t *testing.T) {
	var b strings.Builder
	opts := DefaultSummaryOptions()
	opts.ShowKeys = true
	Summary(&b, sampleReport(), opts)
	out := b.String()
	if !strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected key name in output, got: %s", out)
	}
	if !strings.Contains(out, "DB_HOST") {
		t.Errorf("expected key name in output, got: %s", out)
	}
}

func TestSummary_EnvironmentCount(t *testing.T) {
	var b strings.Builder
	Summary(&b, sampleReport(), DefaultSummaryOptions())
	out := b.String()
	if !strings.Contains(out, "3 environment(s)") {
		t.Errorf("expected environment count, got: %s", out)
	}
}

func TestSummary_NoKeysWhenDisabled(t *testing.T) {
	var b strings.Builder
	opts := DefaultSummaryOptions()
	opts.ShowKeys = false
	Summary(&b, sampleReport(), opts)
	out := b.String()
	if strings.Contains(out, "SECRET_KEY") {
		t.Errorf("expected no key names when ShowKeys=false, got: %s", out)
	}
}
