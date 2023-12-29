package http2test

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
)

const (
	comment       = "//"
	mapDelimiter  = ":"
	hostName      = "host"
	pairDelimiter = ","
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func ReadRequest(u *url.URL) (*http.Request, error) {
	if u == nil {
		return nil, errors.New("error: URL is nil")
	}
	if !uri.IsFileScheme(u) {
		return nil, errors.New(fmt.Sprintf("error: invalid URL scheme : %v", u.Scheme))
	}
	buf, err := os.ReadFile(uri.FileName(u))
	if err != nil {
		return nil, err
	}
	byteReader := bytes.NewReader(buf)
	reader := bufio.NewReader(byteReader)
	req, err1 := http.ReadRequest(reader)
	if err1 != nil {
		return nil, err1
	}
	bytes, err2 := ReadContent(buf)
	if err2 != nil {
		return req, err
	}
	if bytes != nil {
		req.Body = nopCloser{bytes}
	}
	return req, nil
}
