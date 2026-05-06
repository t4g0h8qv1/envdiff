package envdiff

import (
	"testing"
)

func TestDetect_NoDrift(t *testing.T) {
	a := map[string]string{"FOO": "bar", "BAZ": "qux"}
	b := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Detect(a, b, "prod", "staging")
	if r.HasDrift() {
		t.Errorf("expected no drift, got %d items", len(r.Items))
	}
}

func TestDetect_MissingInB(t *testing.T) {
	a := map[string]string{"FOO": "bar", "ONLY_A": "val"}
	b := map[string]string{"FOO": "bar"}
	r := Detect(a, b, "prod", "staging")
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
	if r.Items[0].Key != "ONLY_A" {
		t.Errorf("expected key ONLY_A, got %s", r.Items[0].Key)
	}
	if r.Items[0].Severity != SeverityWarning {
		t.Errorf("expected warning severity")
	}
}

func TestDetect_MissingInA(t *testing.T) {
	a := map[string]string{"FOO": "bar"}
	b := map[string]string{"FOO": "bar", "ONLY_B": "val"}
	r := Detect(a, b, "prod", "staging")
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
	if r.Items[0].Key != "ONLY_B" {
		t.Errorf("expected key ONLY_B, got %s", r.Items[0].Key)
	}
}

func TestDetect_ValueConflict_InfoSeverity(t *testing.T) {
	a := map[string]string{"DB_HOST": "localhost"}
	b := map[string]string{"DB_HOST": "prod.db.internal"}
	r := Detect(a, b, "dev", "prod")
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
	if r.Items[0].Severity != SeverityInfo {
		t.Errorf("expected info severity, got %s", r.Items[0].Severity)
	}
}

func TestDetect_SecretKey_CriticalSeverity(t *testing.T) {
	a := map[string]string{"SECRET_KEY": "abc"}
	b := map[string]string{"SECRET_KEY": "xyz"}
	r := Detect(a, b, "dev", "prod")
	if len(r.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(r.Items))
	}
	if r.Items[0].Severity != SeverityCritical {
		t.Errorf("expected critical severity, got %s", r.Items[0].Severity)
	}
}

func TestBySeverity(t *testing.T) {
	r := &DriftReport{
		Items: []DriftItem{
			{Key: "A", Severity: SeverityWarning},
			{Key: "B", Severity: SeverityCritical},
			{Key: "C", Severity: SeverityWarning},
		},
	}
	warnings := r.BySeverity(SeverityWarning)
	if len(warnings) != 2 {
		t.Errorf("expected 2 warnings, got %d", len(warnings))
	}
	critical := r.BySeverity(SeverityCritical)
	if len(critical) != 1 {
		t.Errorf("expected 1 critical, got %d", len(critical))
	}
}
