package merge

import (
	"strings"
	"testing"
)

func TestMerge_StrategyLeft_NoConflict(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{"C": "3"}

	res, err := Merge(left, right, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["A"] != "1" || res.Env["B"] != "2" || res.Env["C"] != "3" {
		t.Errorf("unexpected env: %v", res.Env)
	}
	if len(res.Warnings) != 0 {
		t.Errorf("expected no warnings, got: %v", res.Warnings)
	}
}

func TestMerge_StrategyLeft_Conflict(t *testing.T) {
	left := map[string]string{"KEY": "left-val"}
	right := map[string]string{"KEY": "right-val"}

	res, err := Merge(left, right, StrategyLeft)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "left-val" {
		t.Errorf("expected left-val, got %q", res.Env["KEY"])
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestMerge_StrategyRight_Conflict(t *testing.T) {
	left := map[string]string{"KEY": "left-val"}
	right := map[string]string{"KEY": "right-val"}

	res, err := Merge(left, right, StrategyRight)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "right-val" {
		t.Errorf("expected right-val, got %q", res.Env["KEY"])
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestMerge_StrategyUnion_KeepsLeft(t *testing.T) {
	left := map[string]string{"KEY": "left-val", "ONLY_LEFT": "yes"}
	right := map[string]string{"KEY": "right-val", "ONLY_RIGHT": "yes"}

	res, err := Merge(left, right, StrategyUnion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Env["KEY"] != "left-val" {
		t.Errorf("union should keep left on conflict, got %q", res.Env["KEY"])
	}
	if res.Env["ONLY_LEFT"] != "yes" || res.Env["ONLY_RIGHT"] != "yes" {
		t.Errorf("union should include all keys: %v", res.Env)
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestMerge_UnknownStrategy(t *testing.T) {
	_, err := Merge(map[string]string{}, map[string]string{}, Strategy("bogus"))
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Errorf("error should mention strategy name, got: %v", err)
	}
}

func TestMerge_EmptyMaps(t *testing.T) {
	res, err := Merge(map[string]string{}, map[string]string{}, StrategyUnion)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Env) != 0 {
		t.Errorf("expected empty env, got %v", res.Env)
	}
}
