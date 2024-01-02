package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/http"
	"reflect"
)

const (
	PingResource          = "ping"
	ContentLength         = "Content-Length"
	ContentType           = "Content-Type"
	writeStatusContentLoc = PkgPath + ":Mux/writeStatusContent"
	bytesLoc              = PkgPath + ":Mux/writeBytes"
)

var m = runtime.NewHandlerMap()

// Handle - add pattern and Http handler mux entry
// TO DO : panic on duplicate handler and pattern combination
func Handle(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	status := m.AddHandler(path, handler)
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
	handler, status := m.GetHandlerFromNID(nid)
	if !status.OK() {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if rsc == PingResource {
		ProcessPing[runtime.Log](w, nid)
		return
	}
	handler(w, r)
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
