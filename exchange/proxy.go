package exchange

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/runtime/runtimetest"
	"net/http"
	"net/url"
)

const (
	proxyLookupLocation = PkgPath + ":lookupHandler"
)

var (
	proxy = runtime.NewProxy()
)

// RegisterHandler - add a map entry for a URI and HttpHandler
func RegisterHandler(uri string, handler func(w http.ResponseWriter, r *http.Request)) runtimetest.Status {
	return proxy.Register(uri, handler)
}

func proxyLookup(u *url.URL) (func(w http.ResponseWriter, r *http.Request), runtimetest.Status) {
	if u == nil {
		return nil, runtimetest.NewStatusError(runtimetest.StatusInvalidArgument, proxyLookupLocation, errors.New("invalid argument: URL is nil"))
	}
	return proxy.Lookup(u.Path)
}
