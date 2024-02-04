package exchange

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"sync"
)

const (
	handlerRegisterLocation  = PkgPath + ":RegisterHandler"
	handlerLookupLocation    = PkgPath + ":Lookup"
	handlerLookupNIDLocation = PkgPath + ":LookupFromNID"
)

var (
	httpProxy = NewProxy()
)

// RegisterHandler - add a map entry for a URI and HttpHandler
func RegisterHandler(uri string, handler func(w http.ResponseWriter, r *http.Request)) *runtime.Status {
	return httpProxy.Register(uri, handler)
}

// Proxy - key value pairs of a URI -> HttpHandler
type Proxy struct {
	m *sync.Map
}

// NewProxy - create a new Proxy
func NewProxy() *Proxy {
	p := new(Proxy)
	p.m = new(sync.Map)
	return p
}

// Register - add an HttpHandler to the proxy
func (p *Proxy) Register(uri string, handler func(w http.ResponseWriter, r *http.Request)) *runtime.Status {
	if len(uri) == 0 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, handlerRegisterLocation, errors.New("invalid argument: path is empty"))
	}
	nid, _, ok := UprootUrn(uri)
	if !ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	if handler == nil {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler is nil: [%v]", uri)))
	}
	_, ok1 := p.m.Load(nid)
	if ok1 {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, handlerRegisterLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler already exists: [%v]", uri)))
	}
	p.m.Store(nid, handler)
	return runtime.StatusOK()
}

// Lookup - get an HttpHandler from the proxy, using a URI as the key
func (p *Proxy) Lookup(uri string) (func(w http.ResponseWriter, r *http.Request), *runtime.Status) {
	nid, _, ok := UprootUrn(uri)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, handlerLookupLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	return p.LookupByNID(nid)
}

// LookupByNID - get an HttpHandler from the proxy, using an NID as a key
func (p *Proxy) LookupByNID(nid string) (func(w http.ResponseWriter, r *http.Request), *runtime.Status) {
	v, ok := p.m.Load(nid)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, handlerLookupNIDLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler does not exist: [%v]", nid)))
	}
	if handler, ok1 := v.(func(w http.ResponseWriter, r *http.Request)); ok1 {
		return handler, runtime.StatusOK()
	}
	return nil, runtime.NewStatus(runtime.StatusInvalidContent)
}
