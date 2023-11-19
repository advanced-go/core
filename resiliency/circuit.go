package resiliency

import (
	"github.com/advanced-go/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

const (
	maxLimit = rate.Limit(100)
)

var cbLocation = PkgUri + "/StatusCircuitBreaker"

// StatusSelectFn - typedef for a function that determines when to select a status
type StatusSelectFn func(status runtime.Status) bool

// StatusCircuitBreaker - Circuit breaker functionality based on a runtime.Status. Configuration provides the
// limit and burst for rate limiting, and a function to determine the selection of statuses.
type StatusCircuitBreaker interface {
	Allow(status runtime.Status) bool
	Limit() rate.Limit
	SetLimit(limit rate.Limit)
	Burst() int
	SetBurst(burst int)
}

type circuitConfig struct {
	limiter *rate.Limiter
	fn      StatusSelectFn
}

// Allow - allow the event based on the status
func (c *circuitConfig) Allow(status runtime.Status) bool {
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
func NewStatusCircuitBreaker(limit rate.Limit, burst int, timeout time.Duration, fn StatusSelectFn) (StatusCircuitBreaker, error) {
	//if t.Limit <= 0 || t.Burst <= 0 {
	//	return nil, errors.New(fmt.Sprintf("error: rate limit or burst is invalid limit = %v burst = %v", t.Limit, t.Burst))
	//}
	//if t.Limit > maxLimit {
	//	return nil, errors.New(fmt.Sprintf("error: rate limit [%v] is greater than the maximum [%v]", t.Limit, maxLimit))
	//}
	//if t.Select == nil {
	//	return nil, errors.New(fmt.Sprintf("error: status select function in nil"))
	//}
	cb := new(circuitConfig)
	cb.limiter = rate.NewLimiter(limit, burst)
	cb.fn = fn
	return cb, nil
}

// CloneStatusCircuitBreaker - create a clone of a StatusCircuitBreaker
func CloneStatusCircuitBreaker(cb StatusCircuitBreaker) StatusCircuitBreaker {
	clone := new(circuitConfig)
	if cfg, ok := any(cb).(*circuitConfig); ok {
		clone.fn = cfg.fn
		clone.limiter = rate.NewLimiter(cb.Limit(), cb.Burst())
	}
	return clone
}
