package runtime

import (
	"strings"
)

type pkg struct{}

const (
	PkgUri = "github.com/advanced-go/core/runtime"
)

// PathFromUri - return a path from a scheme less uri
func PathFromUri(rawUri string) string {
	i := strings.Index(rawUri, "/")
	if i < 0 {
		return "[uri invalid]"
	}
	return rawUri[i:]
}
