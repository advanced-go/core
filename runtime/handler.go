package runtime

import (
	"errors"
	"fmt"
	uri2 "github.com/advanced-go/core/uri"
	"net/http"
	"sync"
)

const (
	handlerAddLocation = PkgPath + ":Mux/add"
	handlerGetLocation = PkgPath + ":Mux/get"
)

type HandlerMap struct {
	m *sync.Map
}

func NewHandlerMap() *HandlerMap {
	h := new(HandlerMap)
	h.m = new(sync.Map)
	return h
}

func (h *HandlerMap) AddHandler(uri string, handler func(w http.ResponseWriter, r *http.Request)) Status {
	if len(uri) == 0 {
		return NewStatusError(StatusInvalidArgument, handlerAddLocation, errors.New("invalid argument: path is empty"))
	}
	nid, _, ok := uri2.UprootUrn(uri)
	if !ok {
		return NewStatusError(StatusInvalidArgument, handlerAddLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	if handler == nil {
		return NewStatusError(StatusInvalidArgument, handlerAddLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler is nil: [%v]", uri)))
	}
	_, ok1 := h.m.Load(nid)
	if ok1 {
		return NewStatusError(StatusInvalidArgument, handlerAddLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler already exists: [%v]", uri)))
	}
	h.m.Store(nid, handler)
	return StatusOK()
}

func (h *HandlerMap) GetHandler(uri string) (func(w http.ResponseWriter, r *http.Request), Status) {
	nid, _, ok := uri2.UprootUrn(uri)
	if !ok {
		return nil, NewStatusError(StatusInvalidArgument, handlerAddLocation, errors.New(fmt.Sprintf("invalid argument: path is invalid: [%v]", uri)))
	}
	return h.GetHandlerFromNID(nid)
}

func (h *HandlerMap) GetHandlerFromNID(nid string) (func(w http.ResponseWriter, r *http.Request), Status) {
	v, ok := h.m.Load(nid)
	if !ok {
		return nil, NewStatusError(StatusInvalidArgument, handlerGetLocation, errors.New(fmt.Sprintf("invalid argument: HTTP handler does not exist: [%v]", nid)))
	}
	if handler, ok1 := v.(func(w http.ResponseWriter, r *http.Request)); ok1 {
		return handler, StatusOK()
	}
	return nil, NewStatus(StatusInvalidContent)
}
