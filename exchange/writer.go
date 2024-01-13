package exchange

import (
	"bytes"
	"io"
	"net/http"
)

type ResponseWriter struct {
	statusCode int
	header     http.Header
	body       *bytes.Buffer
}

func NewResponseWriter() *ResponseWriter {
	w := new(ResponseWriter)
	w.header = make(http.Header)
	w.body = new(bytes.Buffer)
	return w
}

func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

func (w *ResponseWriter) Header() http.Header {
	return w.header
}

func (w *ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

func (w *ResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func (w *ResponseWriter) Response() *http.Response {
	r := new(http.Response)
	if w.statusCode == 0 {
		r.StatusCode = http.StatusOK
	} else {
		r.StatusCode = w.statusCode
	}
	r.Header = w.header
	r.Body = io.NopCloser(bytes.NewReader(w.body.Bytes()))
	return r
}
