package controller

import (
	"errors"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

const (
	TimeoutFlag = "TO"
)

type Controller struct {
	RouteName string
	Timeout   *Timeout
	Router    *Router
}

func NewController(routeName string, d time.Duration, primary, secondary *Resource) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Timeout = new(Timeout)
	c.Timeout.Duration = d
	c.Router = NewRouter(primary, secondary)
	return c
}

func (c *Controller) Do(do func(r *http.Request) (*http.Response, *runtime.Status), req *http.Request) (resp *http.Response, status *runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	rsc := c.Router.RouteTo()
	duration := c.duration(req)
	traffic := access.InternalTraffic
	flags := ""
	start := time.Now().UTC()
	if rsc.internal {
		req, resp, status = doInternal(duration, rsc.handler, req)
	} else {
		traffic = access.EgressTraffic
		req.URL = rsc.BuildUri(req.URL)
		if req.URL != nil {
			req.Host = req.URL.Host
		}
		if duration <= 0 {
			resp, status = do(req)
		} else {
			resp, status = doEgress(duration, do, req)
		}
	}
	if resp.StatusCode == http.StatusGatewayTimeout {
		flags = TimeoutFlag
	}
	access.Log(traffic, start, time.Since(start), req, resp, c.RouteName, rsc.Name, Milliseconds(duration), flags)
	return
}

func (c *Controller) duration(req *http.Request) time.Duration {
	var duration time.Duration
	if c.Timeout != nil {
		duration = c.Timeout.Duration
	}
	if ct, ok := req.Context().Deadline(); ok {
		duration = time.Until(ct) * -1
	}
	return duration
}
