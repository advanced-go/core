package resiliency

import (
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

// LogController - an interface that manages logging for a runtime.TypeHandlerFn
type LogController interface {
	Apply(r *http.Request, body any) (t any, status *runtime.Status)
}

type logController struct {
	name    string
	handler runtime.TypeHandlerFn
	log     LogFn
}

// NewLogController - create a controller for only logging
func NewLogController(name string, handler runtime.TypeHandlerFn, log LogFn) LogController {
	ctrl := new(logController)
	ctrl.name = name
	ctrl.handler = handler
	ctrl.log = log
	return ctrl
}

// Apply - call the log controller for each request
func (ctrl *logController) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()

	t, status = ctrl.handler(r, body)
	resp := http.Response{StatusCode: status.Code()}
	if ctrl.log != nil {
		ctrl.log(internalTraffic, start, time.Since(start), r, &resp, ctrl.name, 0, "")
	}
	return
}
