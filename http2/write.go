package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

var (
	writeLoc = PkgPath + ":WriteResponse"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Only supports []byte, string, io.Reader, and io.ReaderCloser for content
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, content any, status runtime.Status, headers any) {
	var e E

	if status == nil {
		status = runtime.StatusOK()
	}
	// if status.Content is available, then that takes precedence
	if status.Content() != nil {
		SetHeaders(w, headers)
		w.WriteHeader(status.Http())
		writeStatusContent[E](w, status, writeLoc)
		return
	}
	if content == nil {
		SetHeaders(w, headers)
		w.WriteHeader(status.Http())
		return
	}
	ct := GetContentType(headers)
	buf, status0 := WriteBytes(content, ct)
	if !status0.OK() {
		e.Handle(status0, status.RequestId(), writeLoc)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	SetHeaders(w, headers)
	if len(ct) == 0 {
		w.Header().Set(ContentType, http.DetectContentType(buf))
	}
	w.WriteHeader(status.Http())
	bytes, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, err), "", "")
	}
	if bytes != len(buf) {
		fmt.Printf(fmt.Sprintf("error on ResponseWriter().Write() -> [got:%v] [want:%v]\n", bytes, len(buf)))
	}
	return
}
