package runtime

import "strings"

func TrimSpace(s string) string {
	return strings.TrimRight(strings.TrimLeft(s, " "), " ")
}
