package runtime

import "reflect"

type pkg struct{}

var (
	PkgUrl = ParsePkgUrl(reflect.TypeOf(any(pkg{})).PkgPath())
	PkgUri = PkgUrl.Host + PkgUrl.Path
)
