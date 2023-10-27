package runtime

import (
	"reflect"
	"strings"
)

type pkg struct{}

var (
	PkgUri  = reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath = PathFromUri(PkgUri)
)

// PathFromUri - return a path from a scheme less uri
func PathFromUri(rawUri string) string {
	i := strings.Index(rawUri, "/")
	if i < 0 {
		return "[uri invalid]"
	}
	return rawUri[i:]
}
