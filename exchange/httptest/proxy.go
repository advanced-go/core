package httptest

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	HttpErrorUri   = "proxy://www.httperror.com"
	BodyIOErrorUri = "proxy://www.bodyioerror.com"
)

var bodyIOErrorResponse = NewIOErrorResponse()

// ErrorProxy - test proxy for http and body I/O errors
func ErrorProxy(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request or URL is nil")
	}
	switch req.URL.String() {
	case HttpErrorUri:
		return HttpErrorProxy(req)
	case BodyIOErrorUri:
		return BodyIOErrorProxy(req)
	}
	return nil, errors.New(fmt.Sprintf("invalid request URL: %v", req.URL.String()))
}

// HttpErrorProxy - reusable http error proxy
func HttpErrorProxy(req *http.Request) (*http.Response, error) {
	return nil, http.ErrHijacked
}

// BodyIOErrorProxy - reusable body I/O error proxy
func BodyIOErrorProxy(req *http.Request) (*http.Response, error) {
	return bodyIOErrorResponse, nil
}

func Pattern(req *http.Request) string {
	if req == nil || req.URL == nil {
		return ""
	}
	pattern := req.URL.Scheme + "://" + req.URL.Host
	if req.URL.Path != "" {
		pattern += "/" + req.URL.Path
	}
	return pattern
}
