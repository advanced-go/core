package startup

import (
	"reflect"
)

type pkg struct{}

var (
	pkgUri = reflect.TypeOf(any(pkg{})).PkgPath()
)
