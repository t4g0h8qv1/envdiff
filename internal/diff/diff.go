package diff

// Type represents the kind of difference found between two env files.
type Type int

const (
	Missing  Type = iota // Key exists in one file but not the other
	Conflict             // Key exists in both files but values differ
)

// Result represents a single difference between two env maps.
type Result struct {
	Key   string
	Type  Type
	Left  string // value from the left/first file (empty if missing)
	Right string // value from the right/second file (empty if missing)
}

// Compare takes two env maps and returns a slice of differences.
func Compare(left, right map[string]string) []Result {
	var results []Result

	for k, lv := range left {
		rv, ok := right[k]
		if !ok {
			results = append(results, Result{Key: k, Type: Missing, Left: lv, Right: ""})
		} else if lv != rv {
			results = append(results, Result{Key: k, Type: Conflict, Left: lv, Right: rv})
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			results = append(results, Result{Key: k, Type: Missing, Left: "", Right: rv})
		}
	}

	return results
}
