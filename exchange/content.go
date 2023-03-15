package exchange

import (
	"github.com/go-http-utils/headers"
	"net/http"
)

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(headers.ContentLocation)
	}
	return ""
}
