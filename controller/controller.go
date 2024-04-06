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
	Router    *Router
}

func NewController(routeName string, primary, secondary *Resource) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Router = NewRouter(primary, secondary)
	return c
}

func (c *Controller) Do(do func(r *http.Request) (*http.Response, *runtime.Status), req *http.Request) (resp *http.Response, status *runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	rsc := c.Router.RouteTo()
	duration := rsc.timeout(req)
	traffic := access.InternalTraffic
	flags := ""
	req.URL = rsc.BuildUri(req.URL)
	if req.URL != nil {
		req.Host = req.URL.Host
	}
	start := time.Now().UTC()
	if rsc.internal {
		req, resp, status = doInternal(duration, rsc.handler, req)
	} else {
		traffic = access.EgressTraffic
		if duration <= 0 {
			resp, status = do(req)
		} else {
			resp, status = doEgress(duration, do, req)
		}
	}
	elapsed := time.Since(start)
	c.Router.UpdateStats(resp.StatusCode, rsc)
	if resp.StatusCode == http.StatusGatewayTimeout {
		flags = TimeoutFlag
	}
	access.Log(traffic, start, elapsed, req, resp, c.RouteName, rsc.Name, Milliseconds(duration), flags)
	return
}
