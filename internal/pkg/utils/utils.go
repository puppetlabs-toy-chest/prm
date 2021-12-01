package utils

// contains checks if a string is present in a slice
func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

// finds a string present in a slice
func Find(source []string, match string) []string {
	var matches []string
	if Contains(source, match) {
		matches = append(matches, match)
	}
	return matches
}
