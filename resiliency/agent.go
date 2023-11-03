package resiliency

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

const (
	maxLimit = rate.Limit(100)
)

type StatusAgent interface {
	Run(quit chan struct{}, status chan *runtime.Status)
}

type runArgs struct {
	limit rate.Limit
	burst int
	dur   time.Duration
}

var runTable = []runArgs{
	{limit: rate.Limit(0.02), burst: 1, dur: time.Millisecond * 40000},
	{limit: rate.Limit(0.04), burst: 1, dur: time.Millisecond * 20000},
	{limit: rate.Limit(0.08), burst: 1, dur: time.Millisecond * 10000},
	{limit: rate.Limit(0.15), burst: 1, dur: time.Millisecond * 5000},
	{limit: rate.Limit(0.30), burst: 1, dur: time.Millisecond * 2500},
	{limit: rate.Limit(0.60), burst: 1, dur: time.Millisecond * 1250},
	{limit: rate.Limit(1.2), burst: 2, dur: time.Millisecond * 625},
	{limit: rate.Limit(2.5), burst: 3, dur: time.Millisecond * 312},
	{limit: rate.Limit(5.0), burst: 5, dur: time.Millisecond * 156}, // 6.4
	{limit: rate.Limit(10), burst: 10, dur: time.Millisecond * 75},  // 13
	{limit: rate.Limit(15), burst: 15, dur: time.Millisecond * 40},
	{limit: rate.Limit(20), burst: 20, dur: time.Millisecond * 35},
	{limit: rate.Limit(25), burst: 25, dur: time.Millisecond * 30},
	{limit: rate.Limit(30), burst: 30, dur: time.Millisecond * 20}, // 40
	{limit: rate.Limit(35), burst: 35, dur: time.Millisecond * 20},
	{limit: rate.Limit(40), burst: 40, dur: time.Millisecond * 20}, // 50
	{limit: rate.Limit(45), burst: 45, dur: time.Millisecond * 20},
	{limit: rate.Limit(50), burst: 50, dur: time.Millisecond * 15}, // 66
	{limit: rate.Limit(55), burst: 55, dur: time.Millisecond * 15},
	{limit: rate.Limit(60), burst: 60, dur: time.Millisecond * 10},
	{limit: rate.Limit(65), burst: 65, dur: time.Millisecond * 10}, // 100
	{limit: rate.Limit(70), burst: 70, dur: time.Millisecond * 8},  // 125
	{limit: rate.Limit(75), burst: 75, dur: time.Millisecond * 5},  // 200
	{limit: rate.Limit(80), burst: 80, dur: time.Millisecond * 5},
	{limit: rate.Limit(85), burst: 85, dur: time.Millisecond * 3},
	{limit: rate.Limit(90), burst: 90, dur: time.Millisecond * 3},
	{limit: rate.Limit(95), burst: 95, dur: time.Millisecond * 3},
	{limit: rate.Limit(100), burst: 100, dur: time.Millisecond * 3}, // 333
	//{limit: rate.Limit(125), burst: 125, dur: time.Millisecond * 3, cb: nil}, //
	//{limit: rate.Limit(250), burst: 250, dur: time.Millisecond * 3, cb: nil}, //
}

func createTable() []runArgs {
	var table []runArgs
	for _, arg := range runTable {
		//arg.cb = NewStatusCircuitBreaker(runTable[i].limit, runTable[i].burst, func(status *runtime.Status) bool { return true })
		//arg.tick = time.Tick(runTable[i].dur)
		table = append(table, arg)
	}
	return table
}

type agentConfig struct {
	timeout time.Duration
	ping    PingFn
	cb      StatusCircuitBreaker
	table   []runArgs
}

// NewStatusAgent -
func NewStatusAgent(timeout time.Duration, ping PingFn, cb StatusCircuitBreaker) (*runtime.Status, StatusAgent) {
	if cb.Limit() > maxLimit {
		return runtime.NewStatusError(runtime.StatusInvalidArgument, "Agent", errors.New(fmt.Sprintf("error: circuit breaker limit [%v] is > 100", cb.Limit()))), nil
	}
	a := new(agentConfig)
	a.timeout = timeout
	a.ping = ping
	a.cb = cb
	a.table = []runArgs{}
	for _, arg := range runTable {
		a.table = append(a.table, arg)
	}
	return runtime.NewStatusOK(), a
}

func (cf *agentConfig) Run(quit chan struct{}, status chan *runtime.Status) {
	go run(cf.table, cf.ping, cf.timeout, cf.cb, quit, status)
}

func run(table []runArgs, ping PingFn, timeout time.Duration, cb StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	start := time.Now().UTC()
	limiterTime := time.Now().UTC()
	i := 0
	targetLimit := cb.Limit()
	cb.SetLimit(runTable[i].limit)
	cb.SetBurst(runTable[i].burst)
	tick := time.Tick(runTable[i].dur)

	for {
		select {
		case <-tick:
			//status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("done -> elapsed time: %v", time.Since(start)), false)
			//return
			ps := callPing(nil, ping, timeout)
			// If exceeded the current limit, then update limiter and increase the tick frequency
			if !cb.Allow(ps) {
				fmt.Printf("target = %v limit = %v dur = %v time = %v elapsed time = %v\n", targetLimit, cb.Limit(), table[i].dur, time.Since(limiterTime), time.Since(start))
				// If having achieved the target, then return
				if cb.Limit() >= targetLimit {
					status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("success -> elapsed time: %v", time.Since(start)), false)
					return
				}
				//status <- runtime.NewStatus(runtime.StatusHaveContent).SetContent(
				//	fmt.Sprintf("target = %v limit = %v dur = %v time = %v elapsed time = %v\n", targetLimit, cb.Limit(), table[i].dur, time.Since(limiterTime), time.Since(start)), false)
				//fmt.Printf("target = %v limit = %v dur = %v time = %v elapsed time = %v\n", targetLimit, cb.Limit(), table[i].dur, time.Since(limiterTime), time.Since(start))
				i++
				if i >= len(table) {
					status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("error: reached end o table -> elapsed time: %v", time.Since(start)), false)
					return
				}
				tick = time.Tick(runTable[i].dur)
				cb.SetLimit(runTable[i].limit)
				cb.SetBurst(runTable[i].burst)
				limiterTime = time.Now().UTC()
			}
		default:
		}
		select {
		case <-quit:
			status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("quit -> elapsed time: %v", time.Since(start)), false)
			return
		default:
		}
	}
}
