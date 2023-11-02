package resiliency

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

const (
	UpstreamTimeoutFlag = "UT"
)

type LogHandler func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)

type Controller interface {
	Apply(r *http.Request, body any) (t any, status *runtime.Status)
}

type Timeout struct {
	StatusCode int
	Duration   time.Duration
}

type Threshold struct {
	Limit int // request per second
	Burst int
}

type ControllerConfig struct {
	Name      string
	Threshold Threshold // requests per second
	Timeout   Timeout
}

type controller struct {
	config         ControllerConfig
	inFailover     bool // if true, then call upstream and also start pinging. If pinging succeeds, then failback
	primaryCircuit StatusCircuitBreaker
	pingCircuit    StatusCircuitBreaker
	ping           func(ctx context.Context) *runtime.Status
	primary        runtime.TypeHandlerFn
	secondary      runtime.TypeHandlerFn
	log            func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)
}

func NewController(cfg ControllerConfig, ping Threshold, primary, secondary runtime.TypeHandlerFn, log LogHandler) Controller {
	ctrl := new(controller)
	return ctrl
}

func (ctrl *controller) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	/*
		if ctrl.inFailover {
			t, status = ctrl.failover()
		} else {
			if ctrl.config.Timeout.Duration > 0 {
				//newReq := r.Clone() //.
			}
			t, status = ctrl.primary(r, body)

		}

	*/
	resp := http.Response{StatusCode: status.Code()}
	//Call log
	//
	// convert Duration to int milliseconds
	ms := 0
	ctrl.log("internal", start, time.Since(start), r, &resp, ctrl.config.Name, ms, statusFlags)
	return t, status
}
