package resiliency

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
)

const (
	maxLimit = rate.Limit(100)
)

var cbLocation = PkgUri + "/StatusCircuitBreaker"

// StatusSelect - typedef for a function that determines whether or not to select a status
type StatusSelect func(status *runtime.Status) bool

// StatusCircuitBreaker - Circuit breaker functionality based on a runtime.Status. Configuration provides the
// limit and burst for rate limiting, and a function to determine the selection of statuses.
type StatusCircuitBreaker interface {
	Allow(status *runtime.Status) bool
	Limit() rate.Limit
	SetLimit(limit rate.Limit)
	Burst() int
	SetBurst(burst int)
}

type circuitConfig struct {
	limiter *rate.Limiter
	fn      StatusSelect
}

// Allow - allow the event based on the status
func (c *circuitConfig) Allow(status *runtime.Status) bool {
	if status == nil {
		return false
	}
	if !c.fn(status) {
		return true
	}
	return c.limiter.Allow()
}

// Limit - rate limit for the circuit
func (c *circuitConfig) Limit() rate.Limit {
	return c.limiter.Limit()
}

// SetLimit - reconfigure the rate limit for the circuit
func (c *circuitConfig) SetLimit(limit rate.Limit) {
	c.limiter.SetLimit(limit)
}

// Burst - burst for the circuit. Burst is used to attenuate temporary spikes in traffic so that the limit is applied fairly across
// the time frame
func (c *circuitConfig) Burst() int {
	return c.limiter.Burst()
}

// SetBurst - reconfigure the burst
func (c *circuitConfig) SetBurst(burst int) {
	c.limiter.SetBurst(burst)
}

// NewStatusCircuitBreaker - create a circuit breaker with argument validation
func NewStatusCircuitBreaker(limit rate.Limit, burst int, fn StatusSelect) (error, StatusCircuitBreaker) {
	if limit <= 0 || burst <= 0 {
		return errors.New(fmt.Sprintf("error: rate limit or burst is invalid limit = %v burst = %v", limit, burst)), nil
	}
	if limit > maxLimit {
		return errors.New(fmt.Sprintf("error: rate limit [%v] is greater than the maximum [%v]", limit, maxLimit)), nil
	}
	if fn == nil {
		return errors.New(fmt.Sprintf("error: status select function in nil")), nil
	}
	cb := new(circuitConfig)
	cb.limiter = rate.NewLimiter(limit, burst)
	cb.fn = fn
	return nil, cb
}
