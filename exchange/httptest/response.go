package httptest

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// NewResponse - create a new response from the provided parameters
func NewResponse(httpStatus int, content []byte, kv ...string) *http.Response {
	if len(kv)&1 == 1 {
		kv = append(kv, "dummy header value")
	}
	resp := &http.Response{StatusCode: httpStatus, Header: make(http.Header), Request: nil}
	resp.Body = newReaderCloser(bytes.NewReader(content), nil)
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		resp.Header.Add(key, kv[i+1])
	}
	return resp
}

// NewIOErrorResponse - create a response that contains a body that will generate an I/O error when read
func NewIOErrorResponse() *http.Response {
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Request: nil}
	resp.Body = newReaderCloser(nil, io.ErrUnexpectedEOF)
	return resp
}
