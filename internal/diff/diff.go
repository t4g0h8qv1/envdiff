package diff

// Result holds the comparison result between two env files.
type Result struct {
	// MissingInRight contains keys present in left but absent in right.
	MissingInRight []string
	// MissingInLeft contains keys present in right but absent in left.
	MissingInLeft []string
	// Conflicts contains keys present in both files but with different values.
	Conflicts []Conflict
}

// Conflict represents a key whose value differs between two env files.
type Conflict struct {
	Key        string
	LeftValue  string
	RightValue string
}

// Compare takes two parsed env maps and returns a Result describing their differences.
func Compare(left, right map[string]string) Result {
	result := Result{}

	for key, leftVal := range left {
		rightVal, exists := right[key]
		if !exists {
			result.MissingInRight = append(result.MissingInRight, key)
			continue
		}
		if leftVal != rightVal {
			result.Conflicts = append(result.Conflicts, Conflict{
				Key:        key,
				LeftValue:  leftVal,
				RightValue: rightVal,
			})
		}
	}

	for key := range right {
		if _, exists := left[key]; !exists {
			result.MissingInLeft = append(result.MissingInLeft, key)
		}
	}

	return result
}

// HasDifferences returns true if the Result contains any differences.
func (r Result) HasDifferences() bool {
	return len(r.MissingInRight) > 0 || len(r.MissingInLeft) > 0 || len(r.Conflicts) > 0
}
