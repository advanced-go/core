package host

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	httpHandlerProxy = NewProxy()
	duration         time.Duration
	authHandler      HttpHandlerFunc //func(w http.ResponseWriter,r *http.Request)
	okFunc           = func(code int) bool { return code == http.StatusOK }
)

func SetHostTimeout(d time.Duration) {
	duration = d
}

func SetAuthHandler(h func(w http.ResponseWriter, r *http.Request), ok func(int) bool) {
	if h != nil {
		authHandler = h
		if ok != nil {
			okFunc = ok
		}
	}
}

// RegisterHandler - add a path and Http handler to the proxy
// TO DO : panic on duplicate handler and pattern combination
func RegisterHandler(path string, handler HttpHandlerFunc) error {
	if len(path) == 0 {
		return errors.New("error: path is empty")
	}
	if handler == nil {
		return errors.New(fmt.Sprintf("error: handler for path %v is nil", path))
	}
	h := handler
	if authHandler != nil {
		h = NewConditionalIntermediary(authHandler, handler, okFunc)
	}
	if duration > 0 {
		h = NewHostTimeoutIntermediary(duration, h)
	}
	err := httpHandlerProxy.Register(path, h)
	return err
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
