package httpx

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

var (
	writeLoc     = PkgUri + "/WriteResponse"
	minWriteLoc  = PkgUri + "/WriteMinResponse"
	serializeLoc = PkgUri + "/serializeContent"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Only supports []byte, string, io.Reader, and io.ReaderCloser for T
func WriteResponse[E runtime.ErrorHandler, T any](w http.ResponseWriter, content T, status *runtime.Status, headers []Attr) {
	var e E

	if status == nil {
		status = runtime.NewStatusOK()
	}
	// if status.Content is available, then that takes precedence
	if status.Content() != nil {
		w.WriteHeader(status.Http())
		SetHeaders(w, headers)
		writeStatusContent[E](w, status, writeLoc)
		return
	}
	buf, status0 := serializeContent[T](content)
	if !status0.OK() {
		e.HandleStatus(status0, status.RequestId(), writeLoc)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headers)
	if w.Header().Get(ContentLength) == "" {
		w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	}
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(status.RequestId(), writeLoc, err)
	}
	return
}

// WriteMinResponse - write a http.Response, with status, optional headers and optional status content
func WriteMinResponse[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, headers []Attr) {
	if status == nil {
		status = runtime.NewStatusOK()
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headers)
	writeStatusContent[E](w, status, minWriteLoc)
}
