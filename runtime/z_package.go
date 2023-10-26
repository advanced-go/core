package runtime

import (
	"reflect"
)

type pkg struct{}

var (
	PkgUri = reflect.TypeOf(any(pkg{})).PkgPath()
)
