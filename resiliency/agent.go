package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

type Ping func(ctx context.Context) *runtime.Status

type StatusAgent interface {
	Run(quit chan struct{}, status chan *runtime.Status)
}

type agentConfig struct {
	probe   time.Duration // how often to probe
	timeout time.Duration
	ping    Ping
	cb      StatusCircuitBreaker
}

func NewStatusAgent(probe, timeout time.Duration, ping Ping, cb StatusCircuitBreaker) StatusAgent {
	a := new(agentConfig)
	a.probe = probe
	a.timeout = timeout
	a.ping = ping
	a.cb = cb
	return a
}

func (cf *agentConfig) Run(quit chan struct{}, status chan *runtime.Status) {
	go run(cf.ping, cf.cb, quit, status)
}

func run(ping Ping, cb StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	for {
		select {
		case <-quit:
			status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent("quit = true", false)
			return
		default:
		}
	}
}

// run - quit and status
func runChannels(_ Ping, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	count := 0
	for {
		select {
		case <-quit:
			c := fmt.Sprintf("received quit with default count: %v\n", count)
			status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent(c, false)
			return
		default:
			count++
		}
	}
}

func runTicks(_ Ping, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	count := 0
	d := 10
	tick := time.Tick(time.Duration(d) * time.Millisecond)

	for {
		select {
		case <-tick:
			count++
			// Tick reset
			if (count % 10) == 0 {
				fmt.Printf("tick count %vms: %v\n", d, count)
				d += 10
				tick = time.Tick(time.Duration(d) * time.Millisecond)
				count = 0
			}
		case <-quit:
			c := fmt.Sprintf("received quit with tick count %vms: %v\n", d, count)
			status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent(c, false)
			return
		default:
		}
	}
}

func runPing(p Ping, target StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	var limit = rate.Limit(1)
	burst := 1
	cb := CloneStatusCircuitBreaker(target, limit, burst)
	start := time.Now().UTC()
	e := NewExponentialDuration(0.5, 250)
	tick := time.Tick(e.Eval())

	fmt.Printf("initial -> limit = %v burst = %v ms = %v\n", limit, burst, e.Current())

	for {
		select {
		case <-tick:
			ps := p(nil)
			// If exceeded the current limit, then update limiter and increase the ticks frequency
			if !cb.Allow(ps) {
				// If having achieved the target, then return
				l1 := cb.Limit()
				l2 := target.Limit()
				if l1 >= l2 {
					status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("success -> elapsed time: %v", time.Since(start)), false)
					return
				}
				tick = time.Tick(e.Eval())
				limit += limit
				burst += burst
				cb = CloneStatusCircuitBreaker(cb, limit, burst)
				fmt.Printf("limit = %v target = %v burst = %v ms = %v\n", limit, target.Limit(), burst, e.Current())
			}
		case <-quit:
			//c := fmt.Sprintf("received quit with ping count: %v\n", count)
			status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("quit -> elapsed time: %v", time.Since(start)), false)
			return
		default:
		}
	}
}
