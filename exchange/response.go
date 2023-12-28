package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	contentType          = "Content-Type"
	contentTypeJson      = "application/json"
	contentTypeText      = "text/plain"
	jsonSuffix           = ".json"
	readResponseLocation = PkgPath + ":readResponse"
)

var (
	http11Bytes = []byte("HTTP/1.1")
	http12Bytes = []byte("HTTP/1.2")
	http20Bytes = []byte("HTTP/2.0")
)

// readResponse - read a Http response given a URL
func readResponse(u *url.URL) (*http.Response, runtime.Status) {
	serverErr := &http.Response{StatusCode: http.StatusInternalServerError, Status: "Internal Error"}

	if u == nil {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New("error: URL is nil"))
	}
	if !uri.IsFileScheme(u) {
		return serverErr, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, errors.New(fmt.Sprintf("error: Invalid URL scheme : %v", u.Scheme)))
	}
	path := u.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	buf, err := os.ReadFile(uri.FileName(u))
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") {
			return &http.Response{StatusCode: http.StatusNotFound, Status: "Not Found"}, runtime.NewStatusError(runtime.StatusInvalidArgument, readResponseLocation, err)
		}
		return serverErr, runtime.NewStatusError(runtime.StatusIOError, readResponseLocation, err)
	}
	if isHttpResponseMessage(buf) {
		resp1, err2 := http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), nil)
		if err2 != nil {
			return serverErr, runtime.NewStatusError(runtime.StatusIOError, readResponseLocation, err)
		}
		return resp1, runtime.StatusOK()
	} else {
		resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(bytes.NewReader(buf))}
		if strings.HasSuffix(path, jsonSuffix) {
			resp.Header.Add(contentType, contentTypeJson)
		} else {
			resp.Header.Add(contentType, contentTypeText)
		}
		return resp, runtime.StatusOK()
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
