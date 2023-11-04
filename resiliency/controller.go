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
	upstreamTimeoutFlag = "UT"
	internalTraffic     = "internal"
	defaultAgentTimeout = time.Hour * 1
)

var controlAgentFailoverloc = PkgUri + "/Controller/failover"

// PingFn - typedef for a ping function that returns a status
type PingFn func(ctx context.Context) *runtime.Status

// LogFn - typedef for a function that provides access logging
type LogFn func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, controllerName string, timeout int, statusFlags string)

// Controller - an interface that manages resiliency for a runtime.TypeHandlerFn
type Controller interface {
	Apply(r *http.Request, body any) (t any, status *runtime.Status)
}

// Threshold - rate limiting, timeout, and status select configuration
type Threshold struct {
	Limit   rate.Limit // request per second
	Burst   int
	Timeout time.Duration
	Select  StatusSelectFn
}

// ControllerConfig - user supplied configuration
type ControllerConfig struct {
	Name         string
	AgentTimeout time.Duration
	Primary      Threshold
	Ping         Threshold
}

type controller struct {
	config         ControllerConfig
	failoverStatus *atomic.Bool
	primaryCircuit StatusCircuitBreaker
	agent          StatusAgent
	primary        runtime.TypeHandlerFn
	secondary      runtime.TypeHandlerFn
	log            LogFn
	e              runtime.ErrorHandler
}

// NewController - create a new resiliency controller
func NewController[E runtime.ErrorHandler](cfg ControllerConfig, primary, secondary runtime.TypeHandlerFn, ping PingFn, log LogFn) (Controller, error) {
	var e E
	if primary == nil || secondary == nil {
		return nil, errors.New(fmt.Sprintf("error: primary [nil:%v] or secondary [nil:%v] is nil", primary == nil, secondary == nil))
	}
	if ping == nil || cfg.Ping.Select == nil {
		return nil, errors.New(fmt.Sprintf("error: ping [nil:%v] or ping status select [nil:%v] is nil", ping == nil, cfg.Ping.Select == nil))
	}
	var err0 error

	ctrl := new(controller)
	ctrl.config = cfg
	// primary circuit
	ctrl.primaryCircuit, err0 = NewStatusCircuitBreaker(cfg.Primary)
	if err0 != nil {
		return nil, err0
	}
	// ping circuit
	cb, err := NewStatusCircuitBreaker(cfg.Ping)
	if err != nil {
		return nil, err
	}
	// status agent
	ctrl.agent, err0 = NewStatusAgent(ctrl.config.Ping.Timeout, ping, cb)
	if err0 != nil {
		return nil, err0
	}
	if ctrl.config.AgentTimeout == 0 {
		ctrl.config.AgentTimeout = defaultAgentTimeout
	}
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
			tick := time.Tick(c.config.AgentTimeout)
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
			select {
			case <-tick:
				c.e.HandleStatus(runtime.NewStatus(runtime.StatusDeadlineExceeded), "", controlAgentFailoverloc)
				done <- struct{}{}
				return
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
		t, status = callPrimary(r, body, c.primary, c.config.Primary.Timeout)
		// check the circuit
		if !c.primaryCircuit.Allow(status) {
			c.failover()
		}
	}
	// access logging
	resp := http.Response{StatusCode: status.Code()}
	d := time.Since(start)
	c.log(internalTraffic, start, d, r, &resp, c.config.Name, int(d/time.Millisecond), statusFlags)
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

type bypass struct {
	handler runtime.TypeHandlerFn
}

func NewBypassController(handler runtime.TypeHandlerFn) Controller {
	ctrl := new(bypass)
	ctrl.handler = handler
	return ctrl
}

func (b *bypass) Apply(r *http.Request, body any) (t any, status *runtime.Status) {
	return b.handler(r, body)
}
