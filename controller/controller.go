package controller

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"time"
)

const (
	RateLimitInfValue = 99999
	EgressTraffic     = "egress"

	DefaultControllerName = "*"
	NilControllerName     = "!"
	NilBehaviorName       = "!"
	FromRouteHeaderName   = "from-route"
	RequestIdHeaderName   = "X-REQUEST-ID"
	RateLimitFlag         = "RL"
	UpstreamTimeoutFlag   = "UT"
)

// State - defines enabled state
type State interface {
	IsEnabled() bool
	IsNil() bool
	Enable()
	Disable()
}

// Controller - definition for properties of a controller
type Controller interface {
	Actuator
	Name() string
	Timeout() Timeout
	RateLimiter() RateLimiter
	UpdateHeaders(req *http.Request)
	Log(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string)
	t() *controller
}

type controller struct {
	name        string
	ping        bool
	tbl         *table
	timeout     *timeout
	rateLimiter *rateLimiter
}

func cloneController[T *timeout | *rateLimiter](curr *controller, item T) *controller {
	newC := new(controller)
	*newC = *curr
	switch i := any(item).(type) {
	case *timeout:
		newC.timeout = i
	case *rateLimiter:
		newC.rateLimiter = i
	default:
	}
	return newC
}

func newController(route Route, t *table) (*controller, []error) {
	var errs []error
	var err error
	ctrl := newDefaultController(route.Name)
	ctrl.tbl = t
	if route.Timeout != nil {
		ctrl.timeout = newTimeout(route.Name, t, route.Timeout)
		err = ctrl.timeout.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if route.RateLimiter != nil {
		ctrl.rateLimiter = newRateLimiter(route.Name, t, route.RateLimiter)
		err = ctrl.rateLimiter.validate()
		if err != nil {
			errs = append(errs, err)
		}
	}
	return ctrl, errs
}

func newDefaultController(name string) *controller {
	ctrl := new(controller)
	ctrl.name = name
	ctrl.timeout = nilTimeout
	ctrl.rateLimiter = nilRateLimiter
	return ctrl
}

func (c *controller) t() *controller {
	return c
}

func (c *controller) Name() string {
	return c.name
}

func (c *controller) Timeout() Timeout {
	return c.timeout
}

func (c *controller) RateLimiter() RateLimiter {
	return c.rateLimiter
}

func (c *controller) Signal(values url.Values) error {
	if values == nil {
		return nil
	}
	switch values.Get(BehaviorKey) {
	case TimeoutBehavior:
		return c.Timeout().Signal(values)
		break
	case RateLimitBehavior:
		return c.RateLimiter().Signal(values)
		break
	}
	return errors.New(fmt.Sprintf("invalid argument: behavior [%s] is not supported", values.Get(BehaviorKey)))
}

func (c *controller) UpdateHeaders(req *http.Request) {
	if req == nil || req.Header == nil {
		return
	}
	req.Header.Add(FromRouteHeaderName, c.name)
	if req.Header.Get(RequestIdHeaderName) == "" {
		req.Header.Add(RequestIdHeaderName, uuid.New().String())
	}
}

func (c *controller) Log(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string) {
	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add(RequestIdHeaderName, requestId)

	limit, burst, threshold := rateLimiterState(c.rateLimiter)
	if defaultExtractFn != nil {
		defaultExtractFn(EgressTraffic, start, duration, req, statusCode, c.Name(), timeoutState(c.timeout), limit, burst, threshold, statusFlags)
	}
	defaultLogFn(EgressTraffic, start, duration, req, statusCode, c.Name(), timeoutState(c.timeout), limit, burst, threshold, statusFlags)
}
