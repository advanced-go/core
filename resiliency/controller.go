package resiliency

import (
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

const (
	UpstreamTimeoutFlag = "UT"
)

type Timeout struct {
	StatusCode int
	Duration   time.Duration
}

type Threshold struct {
	Percent  int
	Duration time.Duration
}

type Controller struct {
	Name        string
	InFailover  bool // if true, then call upstream and also start pinging. If pinging succeeds, then failback
	FailoverUri string
	PingUri     string
	Timeout     Timeout
	Threshold   Threshold
	Log         func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)
	Handler     runtime.TypeHandlerFn
}

func (ctrl *Controller) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	if ctrl.InFailover {
		t, status = ctrl.failover()
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

func (ctrl *Controller) failover() (any, *runtime.Status) {
	return nil, nil
}

// SetFailover - manual set/reset failover status
func (ctrl *Controller) SetFailover(status bool) {

}
