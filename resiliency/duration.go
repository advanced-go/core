package resiliency

import (
	"math"
	"time"
)

type ExponentialDuration struct {
	rate    float64
	initial float64
	current float64
	count   int
}

func NewExponentialDuration(rate float64, initial float64) *ExponentialDuration {
	e := new(ExponentialDuration)
	e.rate = rate
	e.initial = initial
	e.current = initial
	return e
}

func (e *ExponentialDuration) Eval() time.Duration {
	if e.count != 0 {
		e.current = e.current * (1 - e.rate)
	}
	e.count++
	return time.Duration(math.Trunc(e.current)) * time.Millisecond
}
