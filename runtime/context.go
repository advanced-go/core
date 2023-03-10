package runtime

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const (
	xRequestIdName = "x-request-id"
)

type contextKey struct {
	name string
}

func (k *contextKey) String() string { return "context value " + k.name }

var (
	requestContextKey = &contextKey{"request-id"}
)

// ContextWithRequestId - creates a new Context with a request id
func ContextWithRequestId(ctx context.Context, requestId string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	} else {
		i := ctx.Value(requestContextKey)
		if i != nil {
			return ctx
		}
	}
	if requestId == "" {
		requestId = uuid.New().String()
	}
	return ContextWithValue(ctx, requestContextKey, requestId)
}

// ContextWithRequest - creates a new Context with a request id from the request headers
func ContextWithRequest(req *http.Request) context.Context {
	if req == nil || req.Header == nil {
		return context.Background()
	}
	if req.Header.Get(xRequestIdName) == "" {
		req.Header.Add(xRequestIdName, uuid.New().String())
	}
	return ContextWithRequestId(req.Context(), req.Header.Get(xRequestIdName))
}

// ContextRequestId - return the requestId from a Context
func ContextRequestId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	i := ctx.Value(requestContextKey)
	if requestId, ok := i.(string); ok {
		return requestId
	}
	return ""
}

// ContextWithValue - create a new context with a value, updating the context if it is an HttpExchange context
func ContextWithValue(ctx context.Context, key any, val any) context.Context {
	if ctx == nil {
		return nil
	}
	if curr, ok := any(ctx).(withValue); ok {
		//ctx2 := curr.t()
		//ctx2.ctx = context.WithValue(ctx2.ctx, key, val)
		//ctx2.ctx = context.WithValue(ctx2.ctx, key, val)
		//return ctx2
		return curr.withValue(key, val)
	}
	return context.WithValue(ctx, key, val)
}
