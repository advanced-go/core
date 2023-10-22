package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// See https://tools.ietf.org/html/rfc6265 for details of each of the fields of the above cookie.

func ReadCookies(req *http.Request) map[string]*http.Cookie {
	if req == nil {
		return nil
	}
	jar := make(map[string]*http.Cookie)
	for _, c := range req.Cookies() {
		jar[c.Name] = c
	}
	return jar
}

func AddHeaders(req *http.Request, header http.Header) {
	if req == nil || header == nil {
		return
	}
	for key, element := range header {
		req.Header.Add(key, createValue(element))
	}
}

func createValue(v []string) string {
	if v == nil {
		return ""
	}
	var value string
	for i, item := range v {
		if i > 0 {
			value += ","
		}
		value += item
	}
	return value
}

func ReadRequest(uri *url.URL) (*http.Request, error) {
	if uri == nil {
		return nil, errors.New("error: Uri is nil")
	}
	if uri.Scheme != "file" {
		return nil, errors.New(fmt.Sprintf("error: Invalid Uri scheme : %v", uri.Scheme))
	}
	buf, err := ReadFile(uri)
	if err != nil {
		return nil, err
	}
	req, err1 := http.ReadRequest(bufio.NewReader(bytes.NewReader(buf)))
	if err1 != nil {
		return nil, err1
	}
	if len(req.Header.Get("Content-Length")) == 0 {
		return req, nil
	}
	cnt, err2 := strconv.Atoi(req.Header.Get("Content-Length"))
	if cnt <= 0 || err2 != nil {
		return nil, errors.New("error: count <= 0")
	}
	return req, nil
}
