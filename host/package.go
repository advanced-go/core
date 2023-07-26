package host

import "reflect"

type pkg struct{}

var PkgUrl = reflect.TypeOf(any(pkg{})).PkgPath()
