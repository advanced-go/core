package strings

import "strings"

func TrimSpace(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, " "), " ")
}

func IsEmpty(s string) bool {
	if s == "" {
		return true
	}
	return strings.TrimLeft(s, " ") == ""
}
