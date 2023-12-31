package access

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
	IngressTraffic  = "ingress"
	failsafeUri     = "https://invalid-uri.com"
)

// Origin - log source location
type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	App        string
	InstanceId string
}

// SetOrigin - initialize the origin
func SetOrigin(o Origin) {
	origin = o
}

// Formatter - log formatting
type Formatter func(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) string

// SetFormatter - override log formatting
func SetFormatter(fn Formatter) {
	if fn != nil {
		formatter = fn
	}
}

// Logger - log function
type Logger func(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName string, routeTo string, threshold int, thresholdFlags string)

// SetLogger - override logging
func SetLogger(fn Logger) {
	if fn != nil {
		logger = fn
	}
}

var (
	internalLogging = false

	origin    = Origin{}
	formatter = defaultFormatter
	logger    = defaultLogger
)

func DisableTestLogger() {
	if runtime.IsDebugEnvironment() {
		logger = nil
	}
}

func EnableTestLogger() {
	if runtime.IsDebugEnvironment() {
		SetLogger(defaultLogger)
	}
}

func DisableInternalLogging() {
	internalLogging = false
}

func EnableInternalLogging() {
	internalLogging = true
}

// Log - access logging
func Log(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	if logger == nil {
		return
	}
	if traffic == InternalTraffic && !internalLogging {
		return
	}
	logger(origin, traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
}

// LogDeferred - deferred accessing logging
func LogDeferred(traffic string, req *http.Request, routeName, routeTo string, threshold int, thresholdFlags string, status *runtime.Status) func() {
	start := time.Now().UTC()
	return func() {
		Log(traffic, start, time.Since(start), req, &http.Response{StatusCode: (*status).Code(), Status: (*status).Description()}, routeName, routeTo, threshold, thresholdFlags)
	}
}

// NewRequest - create a new request
func NewRequest(h http.Header, method, uri string) *http.Request {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		req, err = http.NewRequest(method, failsafeUri, nil)
	}
	req.Header = h
	return req
}
