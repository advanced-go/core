package exchange

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/url"
)

const (
	endpointRegisterLocation = PkgPath + ":RegisterEndpoint"
	endpointGetLocation      = PkgPath + ":getEndpoint"
)

var (
	endpoints = runtime.NewHandlerMap()
)

func RegisterEndpoint(uri string, handler func(w http.ResponseWriter, r *http.Request)) runtime.Status {
	return endpoints.AddHandler(uri, handler)
}

func getEndpoint(u *url.URL) (func(w http.ResponseWriter, r *http.Request), runtime.Status) {
	if u == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, endpointGetLocation, errors.New("invalid argument: URL is nil"))
	}
	return endpoints.GetHandler(u.Path)
}
