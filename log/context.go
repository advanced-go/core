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

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

var (
	accessFnContextKey = &contextKey{"access-log-fn"}
)

// NewAccessContext - creates a new Context with an access log function
func NewAccessContext(ctx context.Context) context.Context {
	if accessFn == nil {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	} else {
		fn := ctx.Value(accessFnContextKey)
		if fn != nil {
			return ctx
		}
	}
	return runtime.ContextWithValue(ctx, accessFnContextKey, accessFn)
}

// AccessFromContext - return the access logger from a context
func AccessFromContext(ctx context.Context) startup.AccessLogFn {
	if ctx == nil {
		return nil
	}
	i := ctx.Value(accessFnContextKey)
	if fn, ok2 := i.(startup.AccessLogFn); ok2 {
		return fn
	}
	return nil
}
