package diff

// Kind describes the type of difference found between two env files.
type Kind string

const (
	// Missing indicates a key exists in one file but not the other.
	Missing Kind = "missing"
	// Conflict indicates a key exists in both files but with different values.
	Conflict Kind = "conflict"
)

// Result represents a single difference between two env maps.
type Result struct {
	Kind       Kind
	Key        string
	LeftValue  string
	RightValue string
	// MissingIn holds the name/label of the file that is missing the key.
	MissingIn string
}

// Compare takes two env maps and optional file labels, returning all differences.
// leftName and rightName are used to populate MissingIn on Missing results.
func Compare(left, right map[string]string, leftName, rightName string) []Result {
	var results []Result

	for key, lv := range left {
		if rv, ok := right[key]; !ok {
			results = append(results, Result{
				Kind:      Missing,
				Key:       key,
				LeftValue: lv,
				MissingIn: rightName,
			})
		} else if lv != rv {
			results = append(results, Result{
				Kind:       Conflict,
				Key:        key,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for key, rv := range right {
		if _, ok := left[key]; !ok {
			results = append(results, Result{
				Kind:       Missing,
				Key:        key,
				RightValue: rv,
				MissingIn:  leftName,
			})
		}
	}

	return results
}
