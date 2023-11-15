package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

var (
	writeLoc    = PkgUri + "/WriteResponse"
	minWriteLoc = PkgUri + "/WriteMinResponse"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Only supports []byte, string, io.Reader, and io.ReaderCloser for T
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, content any, status *runtime.Status, headers any) {
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
	if content == nil {
		w.WriteHeader(status.Http())
		SetHeaders(w, headers)
		return
	}
	buf, rc, status0 := WriteBytes(content, GetContentType(headers))
	if !status0.OK() {
		e.Handle(status0, status.RequestId(), writeLoc)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headers)
	w.Header().Set(ContentType, rc)
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, err), "", "")
	}
	return
}

// writeMinResponse - write a http.Response, with status, optional headers and optional status content
func writeMinResponse[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, headers any) {
	if status == nil {
		status = runtime.NewStatusOK()
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headers)
	writeStatusContent[E](w, status, minWriteLoc)
}
