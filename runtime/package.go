package runtime

import (
	"strings"
)

type pkg struct{}

const (
	PkgPath = "github.com/advanced-go/core/runtime"
)

// PathFromUri2 - return a path from a scheme less uri
func PathFromUri2(rawUri string) string {
	i := strings.Index(rawUri, "/")
	if i < 0 {
		return "[uri invalid]"
	}
	return rawUri[i:]
}
