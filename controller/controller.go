package controller

import "time"

const (
	HostRouteName = "host"
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

func NewHostController(d time.Duration) *Controller {
	c := new(Controller)
	c.RouteName = HostRouteName

	c.Timeout = new(Timeout)
	c.Timeout.Duration = d
	return c
}

func NewController(routeName string, d time.Duration, primeHost, secondHost, livenessPath string) *Controller {
	c := new(Controller)
	c.RouteName = routeName
	c.Timeout = new(Timeout)
	c.Timeout.Duration = d

	c.Router = new(Router)

	return c
}
