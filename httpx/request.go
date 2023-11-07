package httpx

import (
	"github.com/go-ai-agent/core/log"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"net/http"
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

func Clone(req *http.Request) *http.Request {
	if req == nil {
		return nil
	}
	requestId := runtime.GetOrCreateRequestId(req)
	if req.Header.Get(runtime.XRequestId) == "" {
		req.Header.Set(runtime.XRequestId, requestId)
	}
	if fn := log.Access(); fn != nil {
		return req.Clone(log.NewAccessContext(req.Context()))
	}
	return req
}
