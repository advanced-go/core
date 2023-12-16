package access

import (
	"github.com/advanced-go/core/runtime"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type pkg struct{}

const (
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
	IngressTraffic  = "ingress"
)

type Origin struct {
	Region     string
	Zone       string
	SubZone    string
	App        string
	InstanceId string
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

//func GetLogger() Logger {
//	return logger
//}

var (
	internalLogging = false
	formatter       = defaultFormatter
	logger          = defaultLogger
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
func Log(o Origin, traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	if logger == nil {
		return
	}
	if traffic == InternalTraffic && !internalLogging {
		return
	}
	logger(o, traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
}

// LogDeferred - deferred accessing logging
func LogDeferred(o Origin, traffic string, req *http.Request, routeName, routeTo string, threshold int, thresholdFlags string, statusCode func() int) func() {
	start := time.Now().UTC()
	return func() {
		Log(o, traffic, start, time.Since(start), req, &http.Response{StatusCode: statusCode()}, routeName, routeTo, threshold, thresholdFlags)
	}
}

// AddRequestId - function copied from package http2
func AddRequestId(req *http.Request) string {
	if req == nil {
		return ""
	}
	id := req.Header.Get(runtime.XRequestId)
	if len(id) == 0 {
		uid, _ := uuid.NewUUID()
		id = uid.String()
		req.Header.Set(runtime.XRequestId, id)
	}
	return id
}

// NewRequest - create a new request
func NewRequest(h http.Header, method, uri string) *http.Request {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		req, err = http.NewRequest(method, "http://invalid-uri.com", nil)
	}
	req.Header = h
	return req
}

// NewStatusCodeClosure - return a func that will return the status code
func NewStatusCodeClosure(status *runtime.Status) func() int {
	return func() int {
		return (*(status)).Code()
	}
}
