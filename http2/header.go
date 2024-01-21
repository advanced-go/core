package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
	"strings"
)

const (
	ContentTypeJson     = "application/json"
	ContentType         = "Content-Type"
	ContentEncoding     = "Content-Encoding"
	AcceptEncoding      = "Accept-Encoding"
	AcceptEncodingValue = "gzip, br"
	ContentLength       = "Content-Length"
	ContentEncodingGzip = "gzip"
	ContentTypeTextHtml = "text/html"
)

func forwardDefaults(dest http.Header, src http.Header) http.Header {
	if dest == nil {
		dest = make(http.Header)
	}
	if src == nil {
		return dest
	}
	// TO DO : add other default headers
	dest.Set(runtime.XRequestId, src.Get(runtime.XRequestId))
	dest.Set(runtime.XRelatesTo, src.Get(runtime.XRelatesTo))
	return dest
}

// Forward - forward headers
func Forward(dest http.Header, src http.Header, names ...string) http.Header {
	dest = forwardDefaults(dest, src)
	if src == nil {
		return dest
	}
	for _, name := range names {
		dest.Set(name, src.Get(name))
	}
	return dest
}

// GetContentType - return the content type header value
func GetContentType(headers any) string {
	if pairs, ok := headers.([]Attr); ok {
		for _, pair := range pairs {
			if pair.Key == ContentType {
				return pair.Val
			}
		}
		return ""
	}
	if h, ok := headers.(http.Header); ok {
		return h.Get(ContentType)
	}
	return ""
}

// SetHeaders - set the headers into an HTTP response writer
func SetHeaders(w http.ResponseWriter, headers any) {
	if headers == nil {
		return
	}
	if pairs, ok := headers.([]Attr); ok {
		for _, pair := range pairs {
			w.Header().Set(strings.ToLower(pair.Key), pair.Val)
		}
		return
	}
	if h, ok := headers.(http.Header); ok {
		for k, v := range h {
			if len(v) > 0 {
				w.Header().Set(strings.ToLower(k), v[0])
			}
		}
	}
}
