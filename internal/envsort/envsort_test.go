package envsort_test

import (
	"testing"

	"github.com/yourusername/envdiff/internal/envsort"
)

func TestSortedKeys_ReturnsAlphaOrder(t *testing.T) {
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	keys := envsort.SortedKeys(env)
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	if len(keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(keys))
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("index %d: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestSortedKeys_EmptyMap(t *testing.T) {
	keys := envsort.SortedKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestNormalize_UppercasesKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "App_Port": "8080"}
	norm := envsort.Normalize(env)
	if v, ok := norm["DB_HOST"]; !ok || v != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", norm)
	}
	if v, ok := norm["APP_PORT"]; !ok || v != "8080" {
		t.Errorf("expected APP_PORT=8080, got %v", norm)
	}
}

func TestSortedPairs_ReturnsSortedKeyValueStrings(t *testing.T) {
	env := map[string]string{"Z": "last", "A": "first"}
	pairs := envsort.SortedPairs(env)
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
	if pairs[0] != "A=first" {
		t.Errorf("expected A=first, got %q", pairs[0])
	}
	if pairs[1] != "Z=last" {
		t.Errorf("expected Z=last, got %q", pairs[1])
	}
}

func TestGroupByPrefix_GroupsCorrectly(t *testing.T) {
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_NAME": "envdiff",
		"NOPREFIX": "value",
	}
	groups := envsort.GroupByPrefix(env)

	if len(groups["DB"]) != 2 {
		t.Errorf("expected 2 DB keys, got %v", groups["DB"])
	}
	if len(groups["APP"]) != 1 {
		t.Errorf("expected 1 APP key, got %v", groups["APP"])
	}
	if len(groups[""]) != 1 {
		t.Errorf("expected 1 unprefixed key, got %v", groups[""])
	}
}
