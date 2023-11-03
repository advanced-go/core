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

type PingFn func(ctx context.Context) *runtime.Status

type LogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)

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
	ping           PingFn
	primary        runtime.TypeHandlerFn
	secondary      runtime.TypeHandlerFn
	log            LogFn
}

func NewController(cfg ControllerConfig, ping PingFn, primary, secondary runtime.TypeHandlerFn, log LogFn) Controller {
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

func callPing(ctx context.Context, ping PingFn, timeout time.Duration) *runtime.Status {
	if timeout == 0 {
		return ping(ctx)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	//ctx, cancel := context.WithTimeout(ctx,timeout)
	status := ping(ctx)

	return status
}
