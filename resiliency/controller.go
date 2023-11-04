package resiliency

import (
	"context"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

const (
	UpstreamTimeoutFlag = "UT"
	InternalTraffic     = "internal"
)

var controllerLoc = PkgUri + "/Controller/ping"

// PingFn - typedef for a ping function that returns a status
type PingFn func(ctx context.Context) *runtime.Status

// LogFn - typedef for a function that provides access logging
type LogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)

// Controller - an interface that manages resiliency for a runtime.TypeHandlerFn
type Controller interface {
	Apply(r *http.Request, body any) (t any, status *runtime.Status)
}

// Timeout - timeout configuration
type Timeout struct {
	StatusCode int
	Duration   time.Duration
}

// Threshold - rate limiting configuration
type Threshold struct {
	Limit int // request per second
	Burst int
}

// ControllerConfig - user supplied configuration
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
	agent          StatusAgent
	ping           PingFn
	primary        runtime.TypeHandlerFn
	secondary      runtime.TypeHandlerFn
	log            LogFn
}

// NewController - create a new resiliency controller
func NewController(cfg ControllerConfig, ping PingFn, primary, secondary runtime.TypeHandlerFn, log LogFn) Controller {
	ctrl := new(controller)
	//ctrl.p
	//ctrl.agent = NewStatusAgent(ctrl.config.Timeout.Duration,ping,pingCircuit)c
	ctrl.config = cfg
	ctrl.ping = ping
	ctrl.primary = primary
	ctrl.secondary = secondary
	ctrl.log = log
	return ctrl
}

func (c *controller) failover() bool {
	return false
}

// Apply - call the controller for each request
func (c *controller) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	if c.failover() {
		t, status = c.secondary(r, body)
	} else {
		t, status = callPrimary(r, body, c.primary, c.config.Timeout.Duration)
		// check the circuit
		if !c.primaryCircuit.Allow(status) {

		}
	}
	// access logging
	resp := http.Response{StatusCode: status.Code()}
	d := time.Since(start)
	c.log(InternalTraffic, start, d, r, &resp, c.config.Name, int(d/time.Millisecond), statusFlags)
	return t, status
}

func callPing(ctx context.Context, fn PingFn, timeout time.Duration) *runtime.Status {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return fn(ctx)
	}
	//ctx, cancel := context.WithTimeout(ctx,timeout)
	status := fn(ctx)
	return status
}

func callPrimary(r *http.Request, body any, fn runtime.TypeHandlerFn, timeout time.Duration) (t any, status *runtime.Status) {
	if timeout <= 0 {
		return fn(r, body)
	}
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	//defer cancel()

	r = r.Clone(ctx)
	t, status = fn(r, body)

	//resp, err = w.rt.RoundTrip(req)
	//if w.deadlineExceeded(err) {
	//	resp = &http.Response{Request: req, StatusCode: tc.StatusCode}
	//	err = nil
	//statusFlags = UpstreamTimeoutFlag
	//	cancel()
	//}
	if status.Code() == runtime.StatusDeadlineExceeded {
		cancel()
	}
	return
}
