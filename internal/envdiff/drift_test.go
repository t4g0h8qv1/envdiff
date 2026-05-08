package envdiff

import (
	"testing"
)

func TestDetectMulti_NoDrift(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"HOST": "localhost", "PORT": "8080"},
		"prod": {"HOST": "prod.example.com", "PORT": "443"},
	}
	// All keys present everywhere — conflicts still exist but no missing
	report := DetectMulti(envs)
	for _, e := range report.Entries {
		if e.Status == "missing" {
			t.Errorf("unexpected missing key: %s", e.Key)
		}
	}
}

func TestDetectMulti_MissingKey(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":     {"HOST": "localhost", "DEBUG": "true"},
		"staging": {"HOST": "staging.example.com"},
		"prod":    {"HOST": "prod.example.com"},
	}
	report := DetectMulti(envs)
	found := false
	for _, e := range report.Entries {
		if e.Key == "DEBUG" && e.Status == "missing" {
			found = true
		}
	}
	if !found {
		t.Error("expected DEBUG to be reported as missing")
	}
}

func TestDetectMulti_Conflict(t *testing.T) {
	envs := map[string]map[string]string{
		"dev":  {"LOG_LEVEL": "debug"},
		"prod": {"LOG_LEVEL": "error"},
	}
	report := DetectMulti(envs)
	if len(report.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(report.Entries))
	}
	if report.Entries[0].Status != "conflict" {
		t.Errorf("expected conflict, got %s", report.Entries[0].Status)
	}
}

func TestDetectMulti_EnvironmentsOrdered(t *testing.T) {
	envs := map[string]map[string]string{
		"prod":    {"X": "1"},
		"dev":     {"X": "1"},
		"staging": {"X": "1"},
	}
	report := DetectMulti(envs)
	expected := []string{"dev", "prod", "staging"}
	for i, e := range expected {
		if report.Environments[i] != e {
			t.Errorf("pos %d: expected %s got %s", i, e, report.Environments[i])
		}
	}
}

func TestDriftReport_Summary_NoDrift(t *testing.T) {
	r := DriftReport{Environments: []string{"dev", "prod"}, Entries: nil}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
	if s != "no drift detected across dev, prod" {
		t.Errorf("unexpected summary: %s", s)
	}
}

func TestDriftReport_Summary_WithDrift(t *testing.T) {
	r := DriftReport{
		Environments: []string{"dev", "prod"},
		Entries: []DriftEntry{
			{Key: "A", Status: "missing"},
			{Key: "B", Status: "conflict"},
		},
	}
	s := r.Summary()
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
