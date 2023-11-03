package resiliency

import (
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
)

type StatusSelect func(status *runtime.Status) bool

// StatusCircuitBreaker - Circuit breaker functionality based on a runtime.Status. Configuration provides the
// limit and burst for rate limiting, and a function to determine the selection of statuses.
type StatusCircuitBreaker interface {
	Allow(status *runtime.Status) bool
	Limit() rate.Limit
	Burst() int
	Select() StatusSelect
}

type circuitConfig struct {
	limiter *rate.Limiter
	fn      StatusSelect
}

func (c *circuitConfig) Allow(status *runtime.Status) bool {
	if status == nil {
		return false
	}
	if !c.fn(status) {
		return true
	}
	return c.limiter.Allow()
}

func (c *circuitConfig) Limit() rate.Limit {
	return c.limiter.Limit()
}

func (c *circuitConfig) Burst() int {
	return c.limiter.Burst()
}

func (c *circuitConfig) Select() StatusSelect {
	return c.fn
}

func NewStatusCircuitBreaker(limit rate.Limit, burst int, fn StatusSelect) StatusCircuitBreaker {
	cb := new(circuitConfig)
	cb.limiter = rate.NewLimiter(limit, burst)
	cb.fn = fn
	return cb
}
