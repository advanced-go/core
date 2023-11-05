package log

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

const (
	XAccessLogger   = "x-access-logger"
	InternalTraffic = "internal"
)

// AccessLogFn - typedef for a function that provides access logging
//type AccessLogFn func(traffic string, start time.Time, duration time.Duration, uri, method string, statusCode int, controllerName string, limit rate.Limit, burst int, timeout int, statusFlags string)

// HttpAccessLogFn - typedef for a function that provides access logging
type HttpAccessLogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

var (
	accessLoggerContextKey = &contextKey{"access-logger"}
)

// ContextWithAccessLogger - creates a new Context with an access logger
func ContextWithAccessLogger(ctx context.Context, access HttpAccessLogFn) context.Context {
	if access == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	} else {
		fn := ctx.Value(accessLoggerContextKey)
		if fn != nil {
			return ctx
		}
	}
	return runtime.ContextWithValue(ctx, accessLoggerContextKey, access)
}

// ContextAccessLogger - return the access logger from a context
func ContextAccessLogger(ctx any) HttpAccessLogFn {
	if ctx == nil {
		return nil
	}
	if ctx2, ok := ctx.(context.Context); ok {
		i := ctx2.Value(accessLoggerContextKey)
		if requestId, ok2 := i.(HttpAccessLogFn); ok2 {
			return requestId
		}
	}
	return nil
}
