package envgroup_test

import (
	"testing"

	"github.com/user/envdiff/internal/envgroup"
)

func makeEnv(pairs ...string) map[string]string {
	m := make(map[string]string)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestByPrefix_GroupsCorrectly(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"DB_PORT", "5432",
		"APP_NAME", "envdiff",
		"NOPREFIX", "value",
	)

	groups := envgroup.ByPrefix(env)
	gm := make(map[string]envgroup.Group)
	for _, g := range groups {
		gm[g.Name] = g
	}

	if len(gm["DB"].Keys) != 2 {
		t.Errorf("expected 2 DB keys, got %d", len(gm["DB"].Keys))
	}
	if len(gm["APP"].Keys) != 1 {
		t.Errorf("expected 1 APP key, got %d", len(gm["APP"].Keys))
	}
	if len(gm["_OTHER"].Keys) != 1 {
		t.Errorf("expected 1 _OTHER key, got %d", len(gm["_OTHER"].Keys))
	}
}

func TestByPrefix_EmptyMap(t *testing.T) {
	groups := envgroup.ByPrefix(map[string]string{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups, got %d", len(groups))
	}
}

func TestByPrefix_ResultIsSorted(t *testing.T) {
	env := makeEnv(
		"Z_KEY", "1",
		"A_KEY", "2",
		"M_KEY", "3",
	)
	groups := envgroup.ByPrefix(env)
	names := make([]string, len(groups))
	for i, g := range groups {
		names[i] = g.Name
	}
	if names[0] != "A" || names[1] != "M" || names[2] != "Z" {
		t.Errorf("groups not sorted: %v", names)
	}
}

func TestByCategories_MatchesPrefixes(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "localhost",
		"REDIS_URL", "redis://",
		"APP_ENV", "production",
		"UNKNOWN", "x",
	)
	cats := map[string][]string{
		"database": {"DB_", "REDIS_"},
		"app":      {"APP_"},
	}
	groups := envgroup.ByCategories(env, cats)
	gm := make(map[string]envgroup.Group)
	for _, g := range groups {
		gm[g.Name] = g
	}

	if len(gm["database"].Keys) != 2 {
		t.Errorf("expected 2 database keys, got %d", len(gm["database"].Keys))
	}
	if len(gm["app"].Keys) != 1 {
		t.Errorf("expected 1 app key, got %d", len(gm["app"].Keys))
	}
	if len(gm["_OTHER"].Keys) != 1 {
		t.Errorf("expected 1 _OTHER key, got %d", len(gm["_OTHER"].Keys))
	}
}

func TestSummary_CountsKeys(t *testing.T) {
	env := makeEnv(
		"DB_HOST", "h",
		"DB_PORT", "p",
		"APP_NAME", "n",
	)
	groups := envgroup.ByPrefix(env)
	summary := envgroup.Summary(groups)

	if summary["DB"] != 2 {
		t.Errorf("expected DB count 2, got %d", summary["DB"])
	}
	if summary["APP"] != 1 {
		t.Errorf("expected APP count 1, got %d", summary["APP"])
	}
}
