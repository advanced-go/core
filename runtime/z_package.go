package runtime

import (
	"reflect"
)

type pkg struct{}

var (
	pkgUri  = reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath = PathFromUri(pkgUri)
)
