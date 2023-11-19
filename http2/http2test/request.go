package http2test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/io2"
	strings2 "github.com/advanced-go/core/strings"
	"io"
	"net/http"
	"net/url"
	"strings"
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
	buf, err := io2.ReadFile(uri)
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
	host, _ := readHostHeader(buf)
	if len(host) > 0 {
		req.Header.Add("Host", host)
		scheme := "https://"
		if strings.Index(host, "local") > -1 {
			scheme = "http://"
		}
		u, _ := url.Parse(scheme + host + req.URL.String())
		req.URL = u
	}
	return req, nil
}

func readHostHeader(buf []byte) (string, error) {
	count := 0
	//m := make(map[string]string)
	r := bytes.NewReader(buf)
	reader := bufio.NewReader(r)
	var line string
	var err error
	for {
		line, err = reader.ReadString('\n')
		count++
		if count == 1 {
			continue
		}
		if isEmpty(line) {
			break
		}
		//k := parseLine(line)
		k, v, err0 := strings2.ParseMapLine(line)
		if err0 != nil {
			return "", err0
		}
		if len(k) > 0 && strings.ToLower(k) == hostName {
			return v, nil
		}
		if err == io.EOF {
			break
		} else {
			if err != nil {
				break
			}
		}
	}
	return "", nil
}

func isEmpty(line string) bool {
	return len(line) == 0 || line == "" || line == "\r\n" || line == "\n"
}
