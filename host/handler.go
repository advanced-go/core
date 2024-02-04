package host

import (
	"net/http"
)

var httpHandlerProxy = NewProxy()

// RegisterHandler - add a path and Http handler to the proxy
// TO DO : panic on duplicate handler and pattern combination
func RegisterHandler(path string, handler ServeHTTPFunc) {
	err := httpHandlerProxy.Register(path, handler)
	if err != nil {
		panic(err)
	}
}

// HttpHandler - handler for messaging
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil || w == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	nid, _, ok := UprootUrn(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler := httpHandlerProxy.LookupByNID(nid)
	if handler == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler(w, r)
}
