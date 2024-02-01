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
