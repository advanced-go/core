package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
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

func runPing(p Ping, cb StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	start := time.Now().UTC()
	count := 0
	tick := time.Tick(100 * time.Millisecond)

	for {
		select {
		case <-tick:
			ps := p(nil)
			if !cb.Allow(ps) {
				status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("elapsed time: %v", time.Since(start)), false)
				return
			}
			count++
		case <-quit:
			c := fmt.Sprintf("received quit with ping count: %v\n", count)
			status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent(c, false)
			return
		default:
		}
	}
}
