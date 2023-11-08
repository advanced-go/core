package log

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
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

var (
	PkgUri   = reflect.TypeOf(any(pkg{})).PkgPath()
	accessFn startup.AccessLogFn
)

func Access() startup.AccessLogFn {
	return accessFn
}

func SetAccess(fn startup.AccessLogFn) {
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
	if accessFn != nil {
		defaultLogFn(EgressTraffic, start, duration, req, resp, threshold, statusFlags)
	}
}

func IngressAccess(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	if accessFn != nil {
		defaultLogFn(IngressTraffic, start, duration, req, resp, threshold, statusFlags)
	}
}

func InternalAccess(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	if accessFn != nil {
		defaultLogFn(InternalTraffic, start, duration, req, resp, threshold, statusFlags)
	}
}

// AnyAccess - needed for packages that have optional logging when core logging is not configured.
func AnyAccess(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	if accessFn != nil {
		defaultLogFn(traffic, start, duration, req, resp, threshold, statusFlags)
	}
}
