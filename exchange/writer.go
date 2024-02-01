package exchange

import (
	"bytes"
	"io"
	"net/http"
)

// ResponseWriter - write a response
type ResponseWriter struct {
	statusCode int
	header     http.Header
	body       *bytes.Buffer
}

// NewResponseWriter - create a new response writer
func NewResponseWriter() *ResponseWriter {
	w := new(ResponseWriter)
	w.header = make(http.Header)
	w.body = new(bytes.Buffer)
	return w
}

// SetStatusCode - return the response status code
func (w *ResponseWriter) SetStatusCode(code int) {
	w.statusCode = code
}

// StatusCode - return the response status code
func (w *ResponseWriter) StatusCode() int {
	return w.statusCode
}

// Header - return the response http.Header
func (w *ResponseWriter) Header() http.Header {
	return w.header
}

// Body - return the response body
func (w *ResponseWriter) Body() []byte {
	return w.body.Bytes()
}

// Write - write the response body
func (w *ResponseWriter) Write(p []byte) (int, error) {
	return w.body.Write(p)
}

// WriteHeader - write the response status code
func (w *ResponseWriter) WriteHeader(statusCode int) {
	if w.statusCode == 0 {
		w.statusCode = statusCode
	}
}

// Response - return the response
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
