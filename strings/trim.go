package strings

import "strings"

// TrimSpace - trim left and right space
func TrimSpace(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, " "), " ")
}

// IsEmpty - determine if a string is empty, and after trimming leading spaces
func IsEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimLeft(s, " ") == ""
}
