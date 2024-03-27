package controller

import "time"

type Timeout struct {
	DurationS string `json:"duration"`
	Duration  time.Duration
}
