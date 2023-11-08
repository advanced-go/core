package log

import (
	"errors"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// Controller - an interface that manages resiliency for a runtime.TypeHandlerFn
type Controller interface {
	Apply(r *http.Request, body any) (t any, status *runtime.Status)
}

type controller struct {
	handler runtime.TypeHandler
}

// NewController - create a new resiliency controller
func NewController(handler runtime.TypeHandler) Controller {
	ctrl := new(controller)
	ctrl.handler = handler
	return ctrl
}

// Apply - call the controller for each request
func (c *controller) Apply(req *http.Request, body any) (any, *runtime.Status) {
	var start = time.Now().UTC()

	if c.handler == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, PkgUri+"/Controller/Apply", errors.New("error: handler function is nil for access logger")).SetRequestId(req.Context())
	}
	t, status := c.handler(req, body)
	InternalAccess(start, time.Since(start), req, &http.Response{StatusCode: status.Code()}, -1, "")
	return t, status
}
