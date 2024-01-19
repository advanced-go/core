package messaging

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	PingResource   = "ping"
	ContentType    = "Content-Type"
	processPingLoc = PkgPath + ":ProcessPing"
)

var proxy = runtime.NewProxy()

// RegisterHandler - add a path and Http handler to the proxy
// TO DO : panic on duplicate handler and pattern combination
func RegisterHandler(path string, handler func(w http.ResponseWriter, r *http.Request)) {
	status := proxy.Register(path, handler)
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
	nid, rsc, ok := uprootUrn(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler, status := proxy.LookupByNID(nid)
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
	var e E
	var buf []byte

	status := Ping[E](nil, nid)
	if status.OK() {
		buf = []byte(fmt.Sprintf("Ping status: %v, resource: %v", status, nid))
	} else {
		buf = []byte(status.FirstError().Error())
	}
	w.Header().Set(ContentType, http.DetectContentType(buf))
	w.WriteHeader(status.Http())
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, processPingLoc, err), "", "")
	}
}
