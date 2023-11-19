package resiliency

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	upstreamTimeoutFlag = "UT"
	internalTraffic     = "internal"
)

type DoHandler func(ctx any, r *http.Request, body any) (any, runtime.Status)

// Controller - an interface that manages resiliency for a runtime.TypeHandlerFn
type Controller interface {
	Apply(r *http.Request, body any) (t any, status runtime.Status)
}

// Threshold - timeout configuration
type Threshold struct {
	Timeout time.Duration
}

type controller struct {
	name      string
	threshold Threshold
	handler   DoHandler
	//log       startup.AccessLogFn
}

// NewController - create a new resiliency controller
func NewController(name string, threshold Threshold, handler DoHandler) Controller {
	//if handler == nil {
	//	return nil, errors.New("error: handler is nil")
	//}
	ctrl := new(controller)
	ctrl.name = name
	ctrl.threshold = threshold
	ctrl.handler = handler
	//ctrl.log = log
	return ctrl
}

func (c *controller) failover() {
	//failoverState := true
	done := make(chan struct{})
	quit := make(chan struct{}, 1)
	status := make(chan runtime.Status, 100)
	go func(chan struct{}, chan runtime.Status) {
		for {
			tick := time.Tick(time.Hour)
			select {
			case st := <-status:
				if st.IsContent() {
					fmt.Printf("test: runTest() -> %v", st.ContentString())
				}
				if st.OK() {
					//failoverState = false
					done <- struct{}{}
					return
				}
			default:
			}
			select {
			case <-tick:
				status <- runtime.NewStatus(runtime.StatusDeadlineExceeded).SetContent("failover", false)
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
func (c *controller) Apply(r *http.Request, body any) (any, runtime.Status) {
	//var start = time.Now().UTC()
	//var statusFlags = ""

	if c.handler == nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, "/Controller/Apply", errors.New(fmt.Sprintf("error: handler function is nil for controller [%v]", c.name))).SetRequestId(r.Context())
	}
	t, status := callHandler(r, body, c.handler, c.threshold.Timeout)
	//resp := http.Response{StatusCode: status.Code()}
	//d := time.Since(start)
	//if c.log != nil {
	//	c.log(internalTraffic, start, d, r, &resp, -1, statusFlags) //c.name, int(c.threshold.Timeout/time.Millisecond), statusFlags)
	//}
	return t, status
}

func callHandler(r *http.Request, body any, fn DoHandler, timeout time.Duration) (t any, status runtime.Status) {
	if timeout <= 0 {
		return fn(nil, r, body)
	}
	ctx, cancel := context.WithTimeout(r.Context(), timeout)
	//defer cancel()

	r = r.Clone(ctx)
	t, status = fn(nil, r, body)

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

/*
type bypass controller

func NewBypassController(name string, handler runtime.TypeHandlerFn, log startup.HttpAccessLogFn) Controller {
	ctrl := new(bypass)
	ctrl.name = name
	ctrl.handler = handler
	ctrl.log = log
	return ctrl
}

func (b *bypass) Apply(r *http.Request, body any) (any, *runtime.Status) {
	var start = time.Now().UTC()
	var statusFlags = ""

	t, status := callHandler(r, body, b.handler, 0)
	resp := http.Response{StatusCode: status.Code()}
	d := time.Since(start)
	if b.log != nil {
		b.log(internalTraffic, start, d, r, &resp, b.name, int(b.threshold.Timeout/time.Millisecond), statusFlags)
	}
	return t, status
}


*/
