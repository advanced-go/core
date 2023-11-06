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
