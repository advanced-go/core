package resiliency

import (
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
)

type Select func(status *runtime.Status) bool

// StatusCircuitBreaker - Circuit breaker functionality based on a runtime.Status. Configuration provides the
// limit and burst for rate limiting, and a function to determine the selection of statuses. The select function
// allows for the creation of an inverted circuit breaker, where exceeding the limit of negative events will break.
// Set the rate limit higher for positive events,and lower for negative events.
type StatusCircuitBreaker interface {
	Allow(status *runtime.Status) bool
	Limit() rate.Limit
}

type circuitConfig struct {
	limiter  *rate.Limiter
	selectFn Select
}

func (c *circuitConfig) Allow(status *runtime.Status) bool {
	if status == nil {
		return false
	}
	if !c.selectFn(status) {
		return true
	}
	return c.limiter.Allow()
}

func (c *circuitConfig) Limit() rate.Limit {
	return c.limiter.Limit()
}

func NewStatusCircuitBreaker(limit rate.Limit, burst int, fn Select) StatusCircuitBreaker {
	cb := new(circuitConfig)
	cb.limiter = rate.NewLimiter(limit, burst)
	cb.selectFn = fn
	return cb
}
