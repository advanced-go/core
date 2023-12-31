package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"reflect"
	"sync"
)

const (
	PingResource          = "ping"
	ContentLength         = "Content-Length"
	ContentType           = "Content-Type"
	muxAddLocation        = PkgPath + ":Mux/add"
	muxGetLocation        = PkgPath + ":Mux/get"
	writeStatusContentLoc = PkgPath + ":Mux/writeStatusContent"
	bytesLoc              = PkgPath + ":Mux/writeBytes"
)

type muxEntry struct {
	path    string
	handler http.HandlerFunc
}

var mux = new(sync.Map)

func add(path string, handler func(w http.ResponseWriter, r *http.Request)) runtime.Status {
	_, ok := mux.Load(path)
	if ok {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, muxAddLocation, errors.New(fmt.Sprintf("invalid argument: HTTP Handler already exists: [%v]", path)))
	}
	entry := new(muxEntry)
	entry.path = path
	entry.handler = handler
	mux.Store(path, entry)
	return runtime.StatusOK()
}

func get(path string) (*muxEntry, runtime.Status) {
	v, ok := mux.Load(path)
	if !ok {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, muxGetLocation, errors.New(fmt.Sprintf("invalid argument: HTTP Handler does not exists: [%v]", path)))
	}
	if entry, ok1 := v.(*muxEntry); ok1 {
		return entry, runtime.StatusOK()
	}
	return nil, runtime.NewStatus(runtime.StatusInvalidContent)
}

// Handle - add pattern and Http handler mux entry
// TO DO : panic on duplicate handler and pattern combination
func Handle(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	status := add(path, handler)
	if !status.OK() {
		panic(status)
	}
}

// HttpHandler - handler for messaging
func HttpHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil || w == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	nid, rsc, ok := uri.UprootUrn(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	entry, status := get(nid)
	if !status.OK() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if rsc == PingResource {
		ProcessPing[runtime.Log](w, nid)
		return
	}
	entry.handler(w, r)
}

func ProcessPing[E runtime.ErrorHandler](w http.ResponseWriter, nid string) {
	status := Ping[E](nil, nid)
	if status.OK() {
		status.SetContent(fmt.Sprintf("Ping status: %v, resource: %v", status, nid), false)
	}
	w.WriteHeader(status.Http())
	writeStatusContent[E](w, status)
}

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status runtime.Status) {
	var e E

	if status.Content() == nil {
		return
	}
	buf, rc, status1 := writeBytes(status.Content())
	if !status1.OK() {
		e.Handle(status, status.RequestId(), writeStatusContentLoc)
		return
	}
	w.Header().Set(ContentType, rc)
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeStatusContentLoc, err), "", "")
	}
}

func writeBytes(content any) ([]byte, string, runtime.Status) {
	var buf []byte

	switch ptr := (content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case error:
		buf = []byte(ptr.Error())
	default:
		return nil, "", runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
	}
	return buf, http.DetectContentType(buf), runtime.StatusOK()
}
