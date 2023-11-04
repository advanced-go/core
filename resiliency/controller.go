package resiliency

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"net/http"
	"sync/atomic"
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

// Threshold - rate limiting and timeout configuration
type Threshold struct {
	Limit    rate.Limit // request per second
	Burst    int
	Duration time.Duration
}

// ControllerConfig - user supplied configuration
type ControllerConfig struct {
	Name    string
	Primary Threshold
	Ping    Threshold
	Agent   Threshold
}

type controller struct {
	config         ControllerConfig
	failoverStatus *atomic.Bool
	primaryCircuit StatusCircuitBreaker
	pingCircuit    StatusCircuitBreaker
	agent          StatusAgent
	ping           PingFn
	primary        runtime.TypeHandlerFn
	secondary      runtime.TypeHandlerFn
	log            LogFn
	e              runtime.ErrorHandler
}

// NewController - create a new resiliency controller
func NewController[E runtime.ErrorHandler](cfg ControllerConfig, primary, secondary runtime.TypeHandlerFn, ping PingFn, statusSelect StatusSelectFn, log LogFn) (Controller, error) {
	var e E
	if primary == nil || secondary == nil {
		return nil, errors.New(fmt.Sprintf("error: primary [nil:%v] or secondary [nil:%v] is nil", primary == nil, secondary == nil))
	}
	if ping == nil || statusSelect == nil {
		return nil, errors.New(fmt.Sprintf("error: ping [nil:%v] or status select [nil:%v] is nil", ping == nil, statusSelect == nil))
	}
	ctrl := new(controller)
	//ctrl.p
	//err,cb := NewStatusCircuitBreaker(0,0,statusSelect)
	//ctrl.agent = NewStatusAgent(ctrl.config.Timeout.Duration,ping,pingCircuit)
	ctrl.config = cfg
	ctrl.ping = ping
	ctrl.primary = primary
	ctrl.secondary = secondary
	ctrl.log = log
	ctrl.failoverStatus = new(atomic.Bool)
	ctrl.failoverStatus.Store(false)
	ctrl.e = e
	return ctrl, nil
}

func (c *controller) failover() {
	if c.failoverStatus.Load() {
		return
	}
	c.failoverStatus.Store(true)
	done := make(chan struct{})
	quit := make(chan struct{}, 1)
	status := make(chan *runtime.Status, 100)
	go func(chan struct{}, chan *runtime.Status) {
		for {
			select {
			case st := <-status:
				if st.IsContent() {
					fmt.Printf("test: runTest() -> %v", st.ContentString())
				}
				if st.OK() {
					c.failoverStatus.Store(false)
					done <- struct{}{}
					return
				}
			default:
			}
		}
	}(done, status)
	<-done
	close(done)
	close(quit)
	close(status)
	return
}

// Apply - call the controller for each request
func (c *controller) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	if c.failoverStatus.Load() {
		t, status = c.secondary(r, body)
	} else {
		t, status = callPrimary(r, body, c.primary, c.config.Primary.Duration)
		// check the circuit
		if !c.primaryCircuit.Allow(status) {
			c.failover()
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
