package model

// The union function returns a slice of strings containing all distinct elements from both slices, preserving order
func union(a, b []string) []string {
	res := make([]string, 0)
	for _, s := range a {
		if !contains(res, s) {
			res = append(res, s)
		}
	}
	for _, s := range b {
		if !contains(res, s) {
			res = append(res, s)
		}
	}
	return res
}

func contains(a []string, val string) bool {
	for _, s := range a {
		if s == val {
			return true
		}
	}
	return false
}
