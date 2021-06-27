package shoebox

import "sort"

func stringSlicesEqual(first []string, second []string) bool {
	if len(first) != len(second) {
		return false
	}

	sort.Strings(first)
	sort.Strings(second)

	for i := range first {
		if first[i] != second[i] {
			return false
		}
	}
	return true
}
