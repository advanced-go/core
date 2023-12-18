package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
	"strings"
)

const (
	//ContentLocation = "Content-Location"
	//ContentTypeText = "text/plain" // charset=utf-8
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
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

/*
func HeaderValue_OLD(name string, r *http.Request) string {
	if r == nil {
		return "invalid-request"
	}
	return r.Header.Get(name)
}

*/

//func GetContentLocation(req *http.Request) string {
//	if req != nil && req.Header != nil {
//		return req.Header.Get(ContentLocation)
//	}
//	return ""
//}

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
		for k, v := range h {
			if k == ContentType {
				if len(v) > 0 {
					return v[0]
				} else {
					return ""
				}
			}
		}
	}
	return ""
}

/*
func CreateHeaders(h http.Header, resp *http.Response, keys ...string) {
	if resp == nil || len(keys) == 0 {
		return
	}
	if keys[0] == "*" {
		keys = []string{}
		for k := range resp.Header {
			keys = append(keys, k)
		}
	}
	if len(keys) > 0 {
		for _, k := range keys {
			if k != "" {
				h.Add(k, resp.Header.Get(k))
			}
		}
	}
}


*/

func SetHeaders(w http.ResponseWriter, headers any) {
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

/*
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

*/
