package resiliency

import (
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

const (
	UpstreamTimeoutFlag = "UT"
)

type TimeoutConfig struct {
	Enabled    bool
	StatusCode int
	Duration   time.Duration
}

type Controller struct {
	Name        string
	InFailover  bool
	FailoverUri string
	PingUri     string
	Timeout     TimeoutConfig
	Log         func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)
	Handler     runtime.TypeHandlerFn
}

func (ctrl *Controller) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	if ctrl.InFailover {
		t, status = failover()
	} else {
		if ctrl.Timeout.Duration > 0 {
			//newReq := r.Clone() //.
		}
		t, status = ctrl.Handler(r, body)

	}
	resp := http.Response{StatusCode: status.Code()}
	//Call log
	//
	// convert Duration to int milliseconds
	ms := 0
	ctrl.Log("egress", start, time.Since(start), r, &resp, ctrl.Name, ms, statusFlags)
	return t, status
}

func failover() (any, *runtime.Status) {
	return nil, nil
}
