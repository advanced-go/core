package access

import (
	"fmt"
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

// LogHandler - access logging handler
type LogHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string)

var (
	handler         LogHandler
	internalLogging = false
)

func GetLogHandler() LogHandler {
	return handler
}

func SetLogHandler(fn LogHandler) {
	if fn != nil {
		handler = fn
	}
}

func DisableTestLogHandler() {
	if runtime.IsDebugEnvironment() {
		handler = nil
	}
}

func EnableTestLogHandler() {
	if runtime.IsDebugEnvironment() {
		SetLogHandler(defaultLogFn)
	}
}

func DisableInternalLogging() {
	internalLogging = false
}

func EnableInternalLogging() {
	internalLogging = true
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	s := fmtLog(traffic, start, duration, req, resp, threshold, thresholdFlags)
	fmt.Printf("%v\n", s)
}

/*
// LogEgress - log egress traffic
func LogEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(EgressTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

// LogIngress - log ingress traffic
func LogIngress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(IngressTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

// LogInternal - log internal package calls
func LogInternal(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(InternalTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

*/

// Log - access logging
func Log(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	if handler == nil {
		return
	}
	if traffic == InternalTraffic && !internalLogging {
		return
	}
	handler(traffic, start, duration, req, resp, threshold, thresholdFlags)
}

// LogDeferred - deferred accessing logging
func LogDeferred(traffic string, req *http.Request, threshold int, thresholdFlags string, statusCode func() int) func() {
	start := time.Now().UTC()
	return func() {
		Log(traffic, start, time.Since(start), req, &http.Response{StatusCode: statusCode()}, threshold, thresholdFlags)
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

// TO DO : Add more header attributes?
func newRequest(h http.Header, method, uri string) *http.Request {
	req, err := http.NewRequest(method, uri, nil)

	if err != nil {
		req, err = http.NewRequest(method, "http://invalid-uri.com", nil)
	}
	requestId := runtime.RequestId(h)
	if len(requestId) > 0 {
		req.Header.Add(runtime.XRequestId, requestId)
	}
	return req
}

// NewStatusCodeClosure - return a func that will return the status code
func NewStatusCodeClosure(status *runtime.Status) func() int {
	return func() int {
		return (*(status)).Code()
	}
}
