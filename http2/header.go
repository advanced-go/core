package http2

import (
	"github.com/go-ai-agent/core/runtime"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

const (
	ContentLocation = "Content-Location"
	ContentTypeText = "text/plain" // charset=utf-8
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
	ContentLength   = "Content-Length"
)

func HeaderValue(name string, r *http.Request) string {
	if r == nil {
		return "invalid-request"
	}
	return r.Header.Get(name)
}

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(ContentLocation)
	}
	return ""
}

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

func AddRequestId(req *http.Request) string {
	if req == nil {
		return ""
	}
	id := req.Header.Get(runtime.XRequestId)
	if len(id) == 0 {
		uid, _ := uuid.NewUUID()
		id = uid.String()
		req.Header.Set(runtime.XRequestId, runtime.GetOrCreateRequestId(req))
	}
	return id
}

/*
func ValidateKVHeaders(kv ...string) error {
	if (len(kv) & 1) == 1 {
		return errors.New("invalid number of kv items: number is odd, missing a value")
	}
	return nil
}

*/
