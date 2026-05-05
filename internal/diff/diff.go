package diff

// Status represents the type of difference found for a key.
type Status string

const (
	StatusMissing  Status = "missing"
	StatusConflict Status = "conflict"
)

// Result holds a single key comparison result between two env maps.
type Result struct {
	Key        string `json:"key"`
	Status     Status `json:"status"`
	LeftValue  string `json:"left_value"`
	RightValue string `json:"right_value"`
}

// Compare takes two env maps and returns a slice of Results describing
// keys that are missing in one side or have conflicting values.
func Compare(left, right map[string]string) []Result {
	var results []Result

	for k, lv := range left {
		if rv, ok := right[k]; !ok {
			results = append(results, Result{
				Key:       k,
				Status:    StatusMissing,
				LeftValue: lv,
			})
		} else if lv != rv {
			results = append(results, Result{
				Key:        k,
				Status:     StatusConflict,
				LeftValue:  lv,
				RightValue: rv,
			})
		}
	}

	for k, rv := range right {
		if _, ok := left[k]; !ok {
			results = append(results, Result{
				Key:        k,
				Status:     StatusMissing,
				RightValue: rv,
			})
		}
	}

	return results
}
