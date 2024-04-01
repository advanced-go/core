package controller

import (
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

type Controller2 struct {
	Name string `json:"name"`
	//Route    string        `json:"route"`
	//Method   string        `json:"method"`
	//Uri      string        `json:"uri"`
	DurationS string `json:"duration"`
	Duration  time.Duration
}

type Controller struct {
	// Identity for access logging route
	RouteName string
	// Selection - how to select this controller given information about the request
	//Path string // package path for selection

	Timeout *Timeout
	Router  *Router
}

func NewTimeoutController(routeName string, d time.Duration) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Timeout = NewTimeout(d)
	return c
}

func NewController(routeName string, d time.Duration) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Timeout = new(Timeout)
	c.Timeout.Duration = d

	c.Router = new(Router)

	return c
}

func (c *Controller) Do(req *http.Request) (*http.Response, *runtime.Status) {
	if req == nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New("invalid argument : request is nil"))
	}
	rsc := c.Router.RouteTo()

	if rsc.internal {

	} else {

	}
	return nil, nil
}
