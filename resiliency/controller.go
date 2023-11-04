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

var controllerLoc = PkgUri + "/Controller/ping"

// PingFn - typedef for a ping function that returns a status
type PingFn func(ctx context.Context) *runtime.Status

// LogFn - typedef for a function that provides access logging
type LogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)

// Controller - an interface that manages resiliency for a configured function of type runtime.TypeHandlerFn
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
	Name    string
	Primary Threshold // requests per second
	Ping    Threshold
	Timeout Timeout
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
	ctrl.config = cfg
	ctrl.ping = ping
	ctrl.primary = primary
	ctrl.secondary = secondary
	ctrl.log = log
	return ctrl
}

// Apply - call the controller for each request
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
	// should not happen, but check for safety
	if ping == nil {
		return runtime.NewStatus(runtime.StatusInvalidArgument).SetContent("error: ping function is nil", false)
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return ping(ctx)
	}
	//ctx, cancel := context.WithTimeout(ctx,timeout)
	status := ping(ctx)
	return status
}
