package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/uri"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	contentType     = "Content-Type"
	contentTypeJson = "application/json"
	contentTypeText = "text/plain"
	jsonSuffix      = ".json"
)

var (
	http11Bytes = []byte("HTTP/1.1")
	http12Bytes = []byte("HTTP/1.2")
	http20Bytes = []byte("HTTP/2.0")
)

// readResponse - read a Http response given a URL
func readResponse(u *url.URL) (*http.Response, error) {
	if u == nil {
		return nil, errors.New("error: URL is nil")
	}
	if u.Scheme != "file" {
		return nil, errors.New(fmt.Sprintf("error: Invalid URL scheme : %v", u.Scheme))
	}
	path := u.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	buf, err := os.ReadFile(uri.FileName(u))
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") {
			return &http.Response{StatusCode: http.StatusNotFound}, err
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, err
	}
	if isHttpResponseMessage(buf) {
		return http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
	} else {
		resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(buf))}
		if strings.HasSuffix(path, jsonSuffix) {
			resp.Header.Add(contentType, contentTypeJson)
		} else {
			resp.Header.Add(contentType, contentTypeText)
		}
		return resp, nil
	}
}

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
