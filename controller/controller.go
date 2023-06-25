package controller

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/time/rate"
	"net/http"
	"net/url"
	"time"
)

const (
	RateLimitInfValue = 99999
	EgressTraffic     = "egress"
	IngressTraffic    = "ingress"
	PingTraffic       = "ping"

	HostControllerName    = "host"
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
	Proxy() Proxy
	UpdateHeaders(req *http.Request)
	LogHttpIngress(start time.Time, duration time.Duration, req *http.Request, statusCode int, written int64, statusFlags string)
	LogHttpEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, retry bool, statusFlags string)
	Log(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string)
	t() *controller
}

type controller struct {
	name        string
	ping        bool
	tbl         *table
	timeout     *timeout
	proxy       *proxy
	rateLimiter *rateLimiter
}

func cloneController[T *timeout | *rateLimiter | *proxy](curr *controller, item T) *controller {
	newC := new(controller)
	*newC = *curr
	switch i := any(item).(type) {
	case *timeout:
		newC.timeout = i
	case *rateLimiter:
		newC.rateLimiter = i
	case *proxy:
		newC.proxy = i
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
	if route.Proxy != nil {
		ctrl.proxy = newProxy(route.Name, t, route.Proxy)
		err = ctrl.proxy.validate()
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
	ctrl.proxy = nilProxy
	return ctrl
}
func (c *controller) validate(egress bool) error {
	if !egress {
		//if c.retry.IsEnabled() {
		//	return errors.New("invalid configuration: Retry is not valid for ingress traffic")
		//	}
		if c.name == HostControllerName {
			if c.timeout.IsEnabled() {
				return errors.New("invalid configuration: Timeout is not valid for host controller")
			}
		} else {
			if c.rateLimiter.IsEnabled() {
				return errors.New("invalid configuration: RateLimiter is not valid for ingress traffic")
			}
		}
	}
	return nil
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

func (c *controller) Proxy() Proxy {
	return c.proxy
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
	case ProxyBehavior:
		return c.Proxy().Signal(values)
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

func (c *controller) LogHttpIngress(start time.Time, duration time.Duration, req *http.Request, statusCode int, written int64, statusFlags string) {
	if c.name == NilControllerName {
		return
	}
	resp := new(http.Response)
	resp.StatusCode = statusCode
	resp.ContentLength = written
	traffic := IngressTraffic
	if c.ping {
		traffic = PingTraffic
	}
	limit, burst, threshold := rateLimiterState(c.rateLimiter)
	proxyValid, proxyThreshold := proxyState(c.proxy)
	if defaultExtractFn != nil {
		defaultExtractFn(traffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
	}
	defaultLogFn(traffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
}

func (c *controller) LogHttpEgress(start time.Time, duration time.Duration, req *http.Request, resp *http.Response, retry bool, statusFlags string) {
	if c.name == NilControllerName {
		return
	}
	var limit rate.Limit
	var burst int
	var threshold string

	limit, burst, threshold = rateLimiterState(c.rateLimiter)
	proxyValid, proxyThreshold := proxyState(c.proxy)
	if defaultExtractFn != nil {
		defaultExtractFn(EgressTraffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
	}
	defaultLogFn(EgressTraffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
}

func (c *controller) Log(start time.Time, duration time.Duration, statusCode int, uri, requestId, method, statusFlags string) {
	req, _ := http.NewRequest(method, uri, nil)
	req.Header.Add(RequestIdHeaderName, requestId)
	resp := new(http.Response)
	resp.StatusCode = statusCode
	limit, burst, threshold := rateLimiterState(c.rateLimiter)
	proxyValid, proxyThreshold := proxyState(c.proxy)
	if defaultExtractFn != nil {
		defaultExtractFn(EgressTraffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
	}
	defaultLogFn(EgressTraffic, start, duration, req, resp, c.Name(), timeoutState(c.timeout), limit, burst, threshold, proxyValid, proxyThreshold, statusFlags)
}
