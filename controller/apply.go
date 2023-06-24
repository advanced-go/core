package controller

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

// Shadowed from : https://grpc.github.io/grpc/core/md_doc_statuscodes.html

const (
	StatusDeadlineExceeded = 4
	StatusRateLimited      = 94
)

type Handler interface {
	Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (fn func(), newCtx context.Context, rateLimited bool)
}

// NilHandler - nil handler
type NilHandler struct{}

// Apply - function to be used to apply a controller
func (e *NilHandler) Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	return nil, nil, false
}

// DefaultHandler - default handler
type DefaultHandler struct{}

// Apply - function to be used to apply a controller
func (e *DefaultHandler) Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	return Apply(ctx, statusCode, uri, requestId, method)
}

// Apply - function to be used to apply a controller
func Apply(ctx context.Context, statusCode func() int, uri, requestId, method string) (func(), context.Context, bool) {
	statusFlags := ""
	limited := false
	start := time.Now()
	newCtx := ctx
	var cancelCtx context.CancelFunc

	ctrl := CtrlTable().LookupUri(uri, method)
	if rlc := ctrl.RateLimiter(); rlc.IsEnabled() && !rlc.Allow() {
		limited = true
		statusFlags = RateLimitFlag
	}
	if !limited {
		if to := ctrl.Timeout(); to.IsEnabled() {
			newCtx, cancelCtx = context.WithTimeout(ctx, to.Duration())
		}
	}
	return func() {
		if cancelCtx != nil {
			cancelCtx()
		}
		code := statusCode()
		if code == StatusDeadlineExceeded {
			statusFlags = UpstreamTimeoutFlag
		}
		ctrl.Log(start, time.Since(start), code, uri, requestId, method, statusFlags)
	}, newCtx, limited
}

func NewStatusCode(status **runtime.Status) func() int {
	return func() int {
		return int((*(status)).Code())
	}
}
