package httpx

import (
	"errors"
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

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(ContentLocation)
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

func SetHeaders(w http.ResponseWriter, kv ...string) error {
	err := ValidateKVHeaders(kv...)
	if err != nil {
		return err
	}
	for i := 0; i < len(kv); i += 2 {
		key := strings.ToLower(kv[i])
		if i+1 >= len(kv) {
			w.Header().Set(key, "")
		} else {
			w.Header().Set(key, kv[i+1])
		}
	}
	return nil
}

func ValidateKVHeaders(kv ...string) error {
	if (len(kv) & 1) == 1 {
		return errors.New("invalid number of kv items: number is odd, missing a value")
	}
	return nil
}
