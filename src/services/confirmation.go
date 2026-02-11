package services

// IsExactMatch checks if the input matches the expected string exactly.
func IsExactMatch(expected, input string) bool {
	return expected == input
}
