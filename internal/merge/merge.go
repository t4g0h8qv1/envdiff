package merge

import (
	"fmt"
	"sort"
)

// Strategy defines how conflicts are resolved during a merge.
type Strategy string

const (
	StrategyLeft  Strategy = "left"  // prefer left file's value on conflict
	StrategyRight Strategy = "right" // prefer right file's value on conflict
	StrategyUnion Strategy = "union" // include all keys, mark conflicts
)

// Result holds the merged key-value map and any warnings generated.
type Result struct {
	Env      map[string]string
	Warnings []string
}

// Merge combines two env maps using the given strategy.
// left and right are maps of key -> value from each respective file.
func Merge(left, right map[string]string, strategy Strategy) (*Result, error) {
	result := &Result{
		Env:      make(map[string]string),
		Warnings: []string{},
	}

	switch strategy {
	case StrategyLeft:
		for k, v := range right {
			result.Env[k] = v
		}
		for k, v := range left {
			if existing, ok := result.Env[k]; ok && existing != v {
				result.Warnings = append(result.Warnings, fmt.Sprintf("conflict on %q: kept left value %q (discarded %q)", k, v, existing))
			}
			result.Env[k] = v
		}
	case StrategyRight:
		for k, v := range left {
			result.Env[k] = v
		}
		for k, v := range right {
			if existing, ok := result.Env[k]; ok && existing != v {
				result.Warnings = append(result.Warnings, fmt.Sprintf("conflict on %q: kept right value %q (discarded %q)", k, v, existing))
			}
			result.Env[k] = v
		}
	case StrategyUnion:
		for k, v := range left {
			result.Env[k] = v
		}
		for k, v := range right {
			if existing, ok := result.Env[k]; ok && existing != v {
				result.Warnings = append(result.Warnings, fmt.Sprintf("conflict on %q: left=%q right=%q (kept left)", k, existing, v))
			} else {
				result.Env[k] = v
			}
		}
	default:
		return nil, fmt.Errorf("unknown merge strategy: %q", strategy)
	}

	sort.Strings(result.Warnings)
	return result, nil
}
