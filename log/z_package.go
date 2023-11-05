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

var (
	PkgUri       = reflect.TypeOf(any(pkg{})).PkgPath()
	accessLogger startup.AccessLogFn
)

func AccessLogger() startup.AccessLogFn {
	return accessLogger
}

func SetAccessLogger(fn startup.AccessLogFn) {
	if fn != nil {
		accessLogger = fn
	}
}
func init() {
	if runtime.IsDebugEnvironment() {
		accessLogger = defaultLogFn
	}
}

var defaultLogFn = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	s := fmtLog(traffic, start, duration, req, resp, threshold, statusFlags)
	fmt.Printf("%v\n", s)
}
