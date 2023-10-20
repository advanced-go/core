package exchange

import (
	"bufio"
	"bytes"
	"github.com/go-ai-agent/core/exchange/httptest"
	"net/http"
	"strings"
)

var http11Bytes = []byte("HTTP/1.1")
var http12Bytes = []byte("HTTP/1.2")
var http20Bytes = []byte("HTTP/2.0")

func isHttpResponseMessage(buf []byte) bool {
	if buf == nil {
		return false
	}
	l := len(http11Bytes)
	if bytes.Equal(buf[0:l], http11Bytes) {
		return true
	}
	l = len(http12Bytes)
	if bytes.Equal(buf[0:l], http12Bytes) {
		return true
	}
	l = len(http20Bytes)
	if bytes.Equal(buf[0:l], http20Bytes) {
		return true
	}
	return false
}

func createFileResponse(req *http.Request) (*http.Response, error) {
	path := req.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	buf, err := ReadFileResource(req.URL)
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") {
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, nil
	}
	if isHttpResponseMessage(buf) {
		return http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), req)
	} else {
		resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: &httptest.ReaderCloser{Reader: bytes.NewReader(buf), Err: nil}}
		if strings.HasSuffix(path, ".json") {
			resp.Header.Add("Content-Type", "application/json")
		} else {
			resp.Header.Add("Content-Type", "text/plain")
		}
		return resp, nil
	}
}
