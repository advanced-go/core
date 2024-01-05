package runtime

import (
	"errors"
	"fmt"
	uri2 "github.com/advanced-go/core/uri"
	"net/http"
	"sync"
)

const (
	handlerRegisterLocation  = PkgPath + ":RegisterHandler"
	handlerLookupLocation    = PkgPath + ":Lookup"
	handlerLookupNIDLocation = PkgPath + ":LookupFromNID"
)

// Proxy - key value pairs of a URI -> HttpHandler
type Proxy struct {
	m *sync.Map
}

// NewProxy - create a new Proxy
func NewProxy() *Proxy {
	h := new(Proxy)
	h.m = new(sync.Map)
	return h
}

// Register - add an HttpHandler to the proxy
func (h *Proxy) Register(uri string, handler func(w http.ResponseWriter, r *http.Request)) Status {
	if len(uri) == 0 {
		return NewStatusError(StatusInvalidArgument, handlerRegisterLocation, errors.New("invalid argument: path is empty"))
	}
	nid, _, ok := uri2.UprootUrn(uri)
	if !ok {
		return NewStatusError(StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	if handler == nil {
		return NewStatusError(StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler is nil: [%v]", uri)))
	}
	_, ok1 := h.m.Load(nid)
	if ok1 {
		return NewStatusError(StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler already exists: [%v]", uri)))
	}
	h.m.Store(nid, handler)
	return StatusOK()
}

// Lookup - get an HttpHandler from the proxy, using a URI as the key
func (h *Proxy) Lookup(uri string) (func(w http.ResponseWriter, r *http.Request), Status) {
	nid, _, ok := uri2.UprootUrn(uri)
	if !ok {
		return nil, NewStatusError(StatusInvalidArgument, handlerLookupLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	return h.LookupByNID(nid)
}

// LookupByNID - get an HttpHandler from the proxy, using an NID as a key
func (h *Proxy) LookupByNID(nid string) (func(w http.ResponseWriter, r *http.Request), Status) {
	v, ok := h.m.Load(nid)
	if !ok {
		return nil, NewStatusError(StatusInvalidArgument, handlerLookupNIDLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler does not exist: [%v]", nid)))
	}
	if handler, ok1 := v.(func(w http.ResponseWriter, r *http.Request)); ok1 {
		return handler, StatusOK()
	}
	return nil, NewStatus(StatusInvalidContent)
}
