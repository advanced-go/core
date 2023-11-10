package log

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"reflect"
	"time"
)

type pkg struct{}

const (
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
	IngressTraffic  = "ingress"
)

// AccessHandler - access logging handler
type AccessHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string)

var (
	PkgUri   = reflect.TypeOf(any(pkg{})).PkgPath()
	accessFn AccessHandler
)

func GetAccessHandler() AccessHandler {
	return accessFn
}

func SetAccessHandler(fn AccessHandler) {
	if fn != nil {
		accessFn = fn
	}
}
func init() {
	if runtime.IsDebugEnvironment() {
		accessFn = defaultLogFn
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
	if accessFn != nil {
		defaultLogFn(traffic, start, duration, req, resp, threshold, statusFlags)
	}
}

// WrapDo - wrap a DoHandler with access logging
func WrapDo(handler runtime.DoHandler) runtime.DoHandler {
	return func(ctx any, req *http.Request, body any) (any, *runtime.Status) {
		var start = time.Now().UTC()

		if handler == nil {
			return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/Wrap", errors.New("error:Do handler function is nil for access log")).SetRequestId(req.Context())
		}
		data, status := handler(ctx, req, body)
		GetAccessHandler()(InternalTraffic, start, time.Since(start), req, &http.Response{StatusCode: status.Code()}, -1, "")
		return data, status
	}
}
