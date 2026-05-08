package envdiff

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func sampleDriftReport() DriftReport {
	return DriftReport{
		Environments: []string{"dev", "prod"},
		Entries: []DriftEntry{
			{
				Key:    "DEBUG",
				Status: "missing",
				Values: map[string]string{"dev": "true", "prod": ""},
			},
			{
				Key:    "LOG_LEVEL",
				Status: "conflict",
				Values: map[string]string{"dev": "debug", "prod": "error"},
			},
		},
	}
}

func TestPrintDrift_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	err := PrintDrift(&buf, sampleDriftReport(), "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "DEBUG") {
		t.Error("expected DEBUG in output")
	}
	if !strings.Contains(out, "LOG_LEVEL") {
		t.Error("expected LOG_LEVEL in output")
	}
	if !strings.Contains(out, "missing") {
		t.Error("expected 'missing' status in output")
	}
	if !strings.Contains(out, "conflict") {
		t.Error("expected 'conflict' status in output")
	}
}

func TestPrintDrift_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	err := PrintDrift(&buf, sampleDriftReport(), "json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out DriftReport
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(out.Entries))
	}
}

func TestPrintDrift_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	report := DriftReport{Environments: []string{"dev", "prod"}, Entries: nil}
	err := PrintDrift(&buf, report, "text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Error("expected 'No drift detected' in output")
	}
}

func TestPrintDrift_DefaultsToText(t *testing.T) {
	var buf bytes.Buffer
	err := PrintDrift(&buf, sampleDriftReport(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.TrimSpace(buf.String()) == "" {
		t.Error("expected non-empty output for default format")
	}
}
