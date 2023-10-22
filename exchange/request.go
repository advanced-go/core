package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	ContentLength = "Content-Length"
)

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

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
	byteReader := bytes.NewReader(buf)
	reader := bufio.NewReader(byteReader)
	req, err1 := http.ReadRequest(reader)
	if err1 != nil {
		return nil, err1
	}
	cnt := contentLength(req)
	if cnt <= 0 {
		return req, nil
	}
	bytes, err2 := readContent(buf)
	if err2 != nil {
		return req, err
	}
	if bytes != nil {
		req.Body = nopCloser{bytes}
	}
	return req, nil
}

func contentLength(req *http.Request) int {
	if req == nil {
		return -1
	}
	s := req.Header.Get(ContentLength)
	if len(s) == 0 {
		return -1
	}
	cnt, err := strconv.Atoi(s)
	if cnt <= 0 || err != nil {
		return -1
	}
	return cnt
}

func readContent(buf []byte) (*bytes.Buffer, error) {
	var content = new(bytes.Buffer)
	var writeTo = false

	reader := bufio.NewReader(bytes.NewReader(buf))
	for {
		line, err := reader.ReadString('\n')
		if len(line) == 2 && !writeTo {
			writeTo = true
			continue
		}
		if writeTo {
			//fmt.Printf("%v", line)
			content.Write([]byte(line))
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				return nil, err
			}
		}
	}
	return content, nil
}
