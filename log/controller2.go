package log

import (
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// Controller2 - an interface that manages resiliency for a runtime.TypeHandlerFn
type Controller2 interface {
	Apply(ctx any, r *http.Request, body any) (any, *runtime.Status)
}

type controller2 struct {
	handler runtime.DoHandlerFn
}

// NewController2 - create a new access logging controller
func NewController2(handler runtime.DoHandlerFn) Controller2 {
	ctrl := new(controller2)
	ctrl.handler = handler
	return ctrl
}

// Apply - call the controller for each request
func (c *controller2) Apply(ctx any, req *http.Request, body any) (any, *runtime.Status) {
	var start = time.Now().UTC()

	if c.handler == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/Controller/Apply", errors.New("error: handler function is nil for access logger")).SetRequestId(req.Context())
	}
	t, status := c.handler(ctx, req, body)
	if fn := AccessFromContext(req.Context()); fn != nil {
		resp := http.Response{StatusCode: status.Code()}
		dur := time.Since(start)
		fn(InternalTraffic, start, dur, req, &resp, -1, "")
	}
	return t, status
}
