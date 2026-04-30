package diff_test

import (
	"testing"

	"github.com/yourorg/envdiff/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	left := map[string]string{"KEY": "value", "PORT": "8080"}
	right := map[string]string{"KEY": "value", "PORT": "8080"}

	result := diff.Compare(left, right)

	if result.HasDifferences() {
		t.Errorf("expected no differences, got %+v", result)
	}
}

func TestCompare_MissingInRight(t *testing.T) {
	left := map[string]string{"KEY": "value", "SECRET": "abc"}
	right := map[string]string{"KEY": "value"}

	result := diff.Compare(left, right)

	if len(result.MissingInRight) != 1 || result.MissingInRight[0] != "SECRET" {
		t.Errorf("expected SECRET missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 0 {
		t.Errorf("expected no keys missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_MissingInLeft(t *testing.T) {
	left := map[string]string{"KEY": "value"}
	right := map[string]string{"KEY": "value", "NEW_KEY": "new"}

	result := diff.Compare(left, right)

	if len(result.MissingInLeft) != 1 || result.MissingInLeft[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY missing in left, got %v", result.MissingInLeft)
	}
}

func TestCompare_Conflicts(t *testing.T) {
	left := map[string]string{"PORT": "8080", "HOST": "localhost"}
	right := map[string]string{"PORT": "9090", "HOST": "localhost"}

	result := diff.Compare(left, right)

	if len(result.Conflicts) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(result.Conflicts))
	}
	c := result.Conflicts[0]
	if c.Key != "PORT" || c.LeftValue != "8080" || c.RightValue != "9090" {
		t.Errorf("unexpected conflict: %+v", c)
	}
}

func TestCompare_Mixed(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"A": "1", "B": "changed", "D": "4"}

	result := diff.Compare(left, right)

	if len(result.MissingInRight) != 1 {
		t.Errorf("expected 1 missing in right, got %v", result.MissingInRight)
	}
	if len(result.MissingInLeft) != 1 {
		t.Errorf("expected 1 missing in left, got %v", result.MissingInLeft)
	}
	if len(result.Conflicts) != 1 {
		t.Errorf("expected 1 conflict, got %v", result.Conflicts)
	}
	if !result.HasDifferences() {
		t.Error("expected HasDifferences to return true")
	}
}
