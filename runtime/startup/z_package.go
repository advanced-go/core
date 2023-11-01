package startup

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"reflect"
)

type pkg struct{}

var (
	PkgUri     = reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath    = runtime.PathFromUri(PkgUri)
	StatusPath = "/startup/status"
)

var StatusRequest = newStatusRequest()

func newStatusRequest() *http.Request {
	req, err := http.NewRequest("", StatusPath, nil)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return req
}
