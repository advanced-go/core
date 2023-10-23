package httpx

import (
	"github.com/go-ai-agent/core/runtime"
	"reflect"
)

type pkg struct{}

var (
	PkgUrl = runtime.ParsePkgUrl(reflect.TypeOf(any(pkg{})).PkgPath())
	PkgUri = PkgUrl.Host + PkgUrl.Path
)
