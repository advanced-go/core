package controller

import "time"

type Controller struct {
	Name string `json:"name"`
	//Route    string        `json:"route"`
	//Method   string        `json:"method"`
	//Uri      string        `json:"uri"`
	DurationS string `json:"duration"`
	Duration  time.Duration
}

type Control2 struct {
	// Identity for access logging route
	RouteName string
	// Selection - how to select this controller given information about the request
	//Path string // package path for selection

	Timeout Timeout
	Router  Router
}
