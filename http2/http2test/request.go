package http2test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/io2/io2test"
	"io"
	"net/http"
	"net/url"
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

func ReadRequest(uri *url.URL) (*http.Request, error) {
	if uri == nil {
		return nil, errors.New("error: Uri is nil")
	}
	if uri.Scheme != "file" {
		return nil, errors.New(fmt.Sprintf("error: Invalid Uri scheme : %v", uri.Scheme))
	}
	buf, err := io2test.ReadFile(uri)
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
