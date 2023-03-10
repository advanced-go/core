package exchange

import (
	"github.com/go-http-utils/headers"
	"net/http"
)

const (
	ContentTypeText = "text/plain" // charset=utf-8
	ContentTypeJson = "application/json"
)

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(headers.ContentLocation)
	}
	return ""
}
