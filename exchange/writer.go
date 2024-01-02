package exchange

import (
	"bytes"
	"io"
	"net/http"
)

type responseWriter struct {
	statusCode int
	header     http.Header
	bytes      *bytes.Buffer
}

func newResponseWriter() *responseWriter {
	w := new(responseWriter)
	w.header = make(http.Header)
	w.bytes = new(bytes.Buffer)
	return w
}

func (w *responseWriter) Header() http.Header {
	return w.header
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return w.bytes.Write(p)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *responseWriter) Result() *http.Response {
	r := new(http.Response)
	if w.statusCode == 0 {
		r.StatusCode = http.StatusOK
	} else {
		r.StatusCode = w.statusCode
	}
	r.Header = w.header
	r.Body = io.NopCloser(bytes.NewReader(w.bytes.Bytes()))
	return r
}
