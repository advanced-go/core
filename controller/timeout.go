package controller

import "time"

type Controller struct {
	Name     string        `json:"name"`
	Route    string        `json:"route"`
	Method   string        `json:"method"`
	Uri      string        `json:"uri"`
	Duration time.Duration `json:"duration"`
}

type Config struct {
	Name     string `json:"name"`
	Route    string `json:"route"`
	Method   string `json:"method"`
	Uri      string `json:"uri"`
	Duration string `json:"duration"`
}
