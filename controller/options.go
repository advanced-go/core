package controller

import (
	"fmt"
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

// UriMatcher - type for control table lookups by uri
type UriMatcher func(uri string, method string) (routeName string, ok bool)

// OutputHandler - type for output handling
type OutputHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, statusCode int, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, statusFlags string)

// SetLogFn - configuration for logging function
func SetLogFn(fn OutputHandler) {
	if fn != nil {
		defaultLogFn = fn
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, statusCode int, routeName string, timeout int, rateLimit rate.Limit, rateBurst int, rateThreshold, statusFlags string) {
	s := FmtLog(traffic, start, duration, req, statusCode, routeName, timeout, rateLimit, rateBurst, rateThreshold, statusFlags)
	fmt.Printf("{%v}\n", s)
}

// SetExtractFn - configuration for connector function
func SetExtractFn(fn OutputHandler) {
	if fn != nil {
		defaultExtractFn = fn
	}
}

var defaultExtractFn OutputHandler
