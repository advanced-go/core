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
	PkgUri          = "github.com/advanced-go/core/access"
	PkgPath         = "/advanced-go/core/access"
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
	IngressTraffic  = "ingress"
)

// LogHandler - access logging handler
type LogHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string)

var (
	handler LogHandler
)

func GetLogHandler() LogHandler {
	return handler
}

func SetLogHandler(fn LogHandler) {
	if fn != nil {
		handler = fn
	}
}

func DisableDebugLogHandler() {
	if runtime.IsDebugEnvironment() {
		handler = nil
	}
}

func EnableDebugLogHandler() {
	if runtime.IsDebugEnvironment() {
		SetLogHandler(defaultLogFn)
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	s := fmtLog(traffic, start, duration, req, resp, threshold, thresholdFlags)
	fmt.Printf("%v\n", s)
}

func LogEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(EgressTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

func LogIngress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(IngressTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

func LogInternal(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	Log(InternalTraffic, start, duration, req, resp, threshold, thresholdFlags)
}

// Log - takes traffic as parameter.
func Log(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, thresholdFlags string) {
	if handler != nil {
		handler(traffic, start, duration, req, resp, threshold, thresholdFlags)
	}
}

// LogDeferred - deferred accessing logging
func LogDeferred(h http.Header, method, uri string, statusCode func() int) func() {
	start := time.Now().UTC()
	req := newRequest(h, method, uri)
	return func() {
		LogInternal(start, time.Since(start), req, &http.Response{StatusCode: statusCode()}, -1, "")
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

func NewStatusCodeClosure(status *runtime.Status) func() int {
	return func() int {
		return (*(status)).Code()
	}
}
