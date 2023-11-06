package log

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/startup"
)

const (
	InternalTraffic = "internal"
	EgressTraffic   = "egress"
)

type key int

var accessFnKey key

// NewAccessContext - creates a new Context with an access log function
func NewAccessContext(ctx context.Context) context.Context {
	if accessFn == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	} else {
		fn := ctx.Value(accessFnKey)
		if fn != nil {
			return ctx
		}
	}
	return runtime.ContextWithValue(ctx, accessFnKey, accessFn)
}

// AccessFromContext - return the access logger from a context
func AccessFromContext(ctx context.Context) startup.AccessLogFn {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(accessFnKey)
	if fn, ok2 := i.(startup.AccessLogFn); ok2 {
		return fn
	}
	return nil
}