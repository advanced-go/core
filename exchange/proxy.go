package exchange

import (
	"errors"
	"github.com/advanced-go/core/runtime"
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
func RegisterHandler(uri string, handler func(w http.ResponseWriter, r *http.Request)) runtime.Status {
	return proxy.Register(uri, handler)
}

func proxyLookup(u *url.URL) (func(w http.ResponseWriter, r *http.Request), runtime.Status) {
	if u == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, proxyLookupLocation, errors.New("invalid argument: URL is nil"))
	}
	return proxy.Lookup(u.Path)
}
