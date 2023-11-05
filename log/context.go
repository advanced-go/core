package log

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
)

const (
	XAccessLogger   = "x-access-logger"
	InternalTraffic = "internal"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

var (
	accessLoggerContextKey = &contextKey{"access-logger"}
)

// ContextWithAccessLogger - creates a new Context with an access logger
func ContextWithAccessLogger(ctx context.Context, access startup.HttpAccessLogFn) context.Context {
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
func ContextAccessLogger(ctx any) startup.HttpAccessLogFn {
	if ctx == nil {
		return nil
	}
	if ctx2, ok := ctx.(context.Context); ok {
		i := ctx2.Value(accessLoggerContextKey)
		if requestId, ok2 := i.(startup.HttpAccessLogFn); ok2 {
			return requestId
		}
	}
	return nil
}
