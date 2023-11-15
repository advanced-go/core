package log2

import (
	"context"
	"errors"
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

// WrapDo - wrap a DoHandler with access logging
func WrapDo(handler runtime.DoHandler) runtime.DoHandler {
	return func(ctx any, req *http.Request, body any) (any, *runtime.Status) {
		var start = time.Now().UTC()

		if handler == nil {
			return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/WrapDo", errors.New("error:Do handler function is nil for access log")).SetRequestId(req.Context())
		}
		data, status := handler(ctx, req, body)
		AnyAccess(InternalTraffic, start, time.Since(start), req, &http.Response{StatusCode: status.Code()}, -1, "")
		return data, status
	}
}

// WrapPost - wrap a PostHandler with access logging
func WrapPost(handler runtime.PostHandler) runtime.PostHandler {
	return func(ctx any, r *http.Request, body any) (any, *runtime.Status) {
		var start = time.Now().UTC()

		//req, _ := http.NewRequest(method, uri, nil)
		if handler == nil {
			return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/WrapPost", errors.New("error:Do handler function is nil for access log")).SetRequestId(r)
		}
		data, status := handler(ctx, r, body)
		AnyAccess(InternalTraffic, start, time.Since(start), r, &http.Response{StatusCode: status.Code()}, -1, "")
		return data, status
	}
}

// WrapHttp - wrap a HttpHandler with access logging
func WrapHttp(handler runtime.HttpHandler) runtime.HttpHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) *runtime.Status {
		var start = time.Now().UTC()

		if handler == nil {
			return runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/WrapHttp", errors.New("error:Http handler function is nil for access log")).SetRequestId(r.Context())
		}
		status := handler(ctx, w, r)
		AnyAccess(InternalTraffic, start, time.Since(start), r, &http.Response{StatusCode: status.Code()}, -1, "")
		return status
	}
}

// AddRequestId - function copied from package httpx
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

// Log - accessing logging for generic function calls
func Log(ctx any, method, uri string, statusCode func() int) func() {
	start := time.Now().UTC()
	req := newRequest(ctx, method, uri)
	return func() {
		InternalAccess(start, time.Since(start), req, &http.Response{StatusCode: statusCode()}, -1, "")
	}
}

func newRequest(ctx any, method, uri string) *http.Request {
	req, err := http.NewRequest(method, uri, nil)

	if err != nil {
		req, err = http.NewRequest(method, "http://invalid-uri.com", nil)
	}
	requestId := runtime.RequestId(ctx)
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
