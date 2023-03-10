package httptest

import (
	"bufio"
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
)

//go:embed resource/*
var fsys embed.FS

// NewResponse - create a new response from the provided parameters
func NewResponse(httpStatus int, content []byte, kv ...string) *http.Response {
	if len(kv)&1 == 1 {
		kv = append(kv, "dummy header value")
	}
	resp := &http.Response{StatusCode: httpStatus, Header: make(http.Header), Request: nil}
	resp.Body = NewReaderCloser(bytes.NewReader(content), nil)
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		resp.Header.Add(key, kv[i+1])
	}
	return resp
}

// ReadResponse - create a response by reading the content from ad embedded file system
func ReadResponse(f fs.FS, name string) (*http.Response, error) {
	buf, err := fs.ReadFile(f, name)
	if err != nil {
		return nil, err
	}
	resp, err0 := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	if err0 != nil {
		return nil, err0
	}
	return resp, nil
}

// ReadResponseTest - only used for calls scoped to the enclosing http package, specifically from do_test.go.
func ReadResponseTest(name string) (*http.Response, error) {
	return ReadResponse(fsys, name)
}

// NewIOErrorResponse - create a respons that contains a body that will generate an I/O error when read
func NewIOErrorResponse() *http.Response {
	resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Request: nil}
	resp.Body = NewReaderCloser(nil, io.ErrUnexpectedEOF)
	return resp
}
