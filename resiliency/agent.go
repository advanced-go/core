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

func runTicks(p Ping, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	count := 0
	tick := time.Tick(10 * time.Millisecond)

	for {
		select {
		case <-tick:
			count++
		case <-quit:
			c := fmt.Sprintf("received quit with ticks count: %v\n", count)
			status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent(c, false)
			return
		default:
			count++
		}
	}
}
