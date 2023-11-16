package log2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type pkg struct{}

const (
	PkgUri          = "github.com/advanced-go/core/log2"
	PkgPath         = "/advanced-go/core/log2"
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
	IngressTraffic  = "ingress"
)

// AccessHandler - access logging handler
type AccessHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string)

var (
	handler AccessHandler
)

func GetAccessHandler() AccessHandler {
	return handler
}

func SetAccessHandler(fn AccessHandler) {
	if fn != nil {
		handler = fn
	}
}

func DisableDebugAccessHandler() {
	if runtime.IsDebugEnvironment() {
		handler = nil
	}
}

func EnableDebugAccessHandler() {
	if runtime.IsDebugEnvironment() {
		SetAccessHandler(defaultLogFn)
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	s := fmtLog(traffic, start, duration, req, resp, threshold, statusFlags)
	fmt.Printf("%v\n", s)
}

func EgressAccess(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	AnyAccess(EgressTraffic, start, duration, req, resp, threshold, statusFlags)
}

func IngressAccess(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	AnyAccess(IngressTraffic, start, duration, req, resp, threshold, statusFlags)
}

func InternalAccess(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	AnyAccess(InternalTraffic, start, duration, req, resp, threshold, statusFlags)
}

// AnyAccess - needed for packages that have optional logging when core logging is not configured.
func AnyAccess(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	if handler != nil {
		handler(traffic, start, duration, req, resp, threshold, statusFlags)
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

// Log - accessing logging
func Log(h http.Header, method, uri string, statusCode func() int) func() {
	start := time.Now().UTC()
	req := newRequest(h, method, uri)
	return func() {
		InternalAccess(start, time.Since(start), req, &http.Response{StatusCode: statusCode()}, -1, "")
	}
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

func NewStatusCodeClosure(status **runtime.Status) func() int {
	return func() int {
		return (*(status)).Code()
	}
}
