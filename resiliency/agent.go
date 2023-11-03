package resiliency

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

type StatusAgent interface {
	Run(quit chan struct{}, status chan *runtime.Status)
}

type runArgs struct {
	limit rate.Limit
	burst int
	dur   time.Duration
	cb    StatusCircuitBreaker
	tick  <-chan time.Time
}

var runTable = []runArgs{
	{limit: rate.Limit(0.15), burst: 1, dur: time.Millisecond * 5000, cb: nil},
	{limit: rate.Limit(0.30), burst: 1, dur: time.Millisecond * 2500, cb: nil},
	{limit: rate.Limit(0.60), burst: 1, dur: time.Millisecond * 1250, cb: nil},
	{limit: rate.Limit(1.2), burst: 2, dur: time.Millisecond * 625, cb: nil},
	{limit: rate.Limit(2.5), burst: 3, dur: time.Millisecond * 312, cb: nil},
	{limit: rate.Limit(5.0), burst: 5, dur: time.Millisecond * 156, cb: nil}, // 6.4
	{limit: rate.Limit(10), burst: 10, dur: time.Millisecond * 75, cb: nil},  // 13
	{limit: rate.Limit(15), burst: 15, dur: time.Millisecond * 40, cb: nil},
	{limit: rate.Limit(20), burst: 20, dur: time.Millisecond * 35, cb: nil},
	{limit: rate.Limit(25), burst: 25, dur: time.Millisecond * 30, cb: nil},
	{limit: rate.Limit(30), burst: 30, dur: time.Millisecond * 20, cb: nil}, // 40
	{limit: rate.Limit(35), burst: 35, dur: time.Millisecond * 20, cb: nil},
	{limit: rate.Limit(40), burst: 40, dur: time.Millisecond * 20, cb: nil}, // 50
	{limit: rate.Limit(45), burst: 45, dur: time.Millisecond * 20, cb: nil},
	{limit: rate.Limit(50), burst: 50, dur: time.Millisecond * 15, cb: nil}, // 66
	{limit: rate.Limit(55), burst: 55, dur: time.Millisecond * 15, cb: nil},
	{limit: rate.Limit(60), burst: 60, dur: time.Millisecond * 10, cb: nil},
	{limit: rate.Limit(65), burst: 65, dur: time.Millisecond * 10, cb: nil}, // 100
	{limit: rate.Limit(70), burst: 70, dur: time.Millisecond * 8, cb: nil},  // 125
	{limit: rate.Limit(75), burst: 75, dur: time.Millisecond * 5, cb: nil},  // 200
	{limit: rate.Limit(80), burst: 80, dur: time.Millisecond * 5, cb: nil},
	{limit: rate.Limit(85), burst: 85, dur: time.Millisecond * 3, cb: nil},
	{limit: rate.Limit(90), burst: 90, dur: time.Millisecond * 3, cb: nil},
	{limit: rate.Limit(95), burst: 95, dur: time.Millisecond * 3, cb: nil},
	{limit: rate.Limit(100), burst: 100, dur: time.Millisecond * 3, cb: nil}, // 333
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

func NewStatusAgent(timeout time.Duration, ping PingFn, cb StatusCircuitBreaker) StatusAgent {
	a := new(agentConfig)
	a.timeout = timeout
	a.ping = ping
	a.cb = cb
	a.table = []runArgs{}
	for _, arg := range runTable {
		//arg.cb = NewStatusCircuitBreaker(runTable[i].limit, runTable[i].burst, cb.Select())
		//arg.tick = time.Tick(runTable[i].dur)
		a.table = append(a.table, arg)
	}
	//cnt := copy(a.table, runTable)
	//fmt.Printf("copy cnt = %v\n", len(a.table))
	return a
}

func (cf *agentConfig) Run(quit chan struct{}, status chan *runtime.Status) {
	go run(cf.ping, cf.table, cf.cb, quit, status)
}

func run(ping PingFn, table []runArgs, cb StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
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
func runChannels(_ PingFn, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
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

func runTicks(_ PingFn, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
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

func runTest(table []runArgs, ping PingFn, timeout time.Duration, cb StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
	start := time.Now().UTC()
	limiterTime := time.Now().UTC()
	i := 0
	targetLimit := cb.Limit()
	cb.SetLimit(runTable[i].limit)
	cb.SetBurst(runTable[i].burst)
	tick := time.Tick(runTable[i].dur)

	//fmt.Printf("target = %v limit = %v dur = %v time = %v elapsed time = %v\n", target.Limit(), cb.Limit(), table[i].dur, time.Since(limiterTime), time.Since(start))
	for {
		select {
		case <-tick:
			ps := callPing(nil, ping, timeout)
			// If exceeded the current limit, then update limiter and increase the tick frequency
			if !cb.Allow(ps) {
				// If having achieved the target, then return
				if cb.Limit() >= targetLimit {
					status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("success -> elapsed time: %v", time.Since(start)), false)
					return
				}
				fmt.Printf("target = %v limit = %v dur = %v time = %v elapsed time = %v\n", targetLimit, cb.Limit(), table[i].dur, time.Since(limiterTime), time.Since(start))
				i++
				if i >= len(table) {
					status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("error: reached end o table -> elapsed time: %v", time.Since(start)), false)
					return
				}
				tick = time.Tick(runTable[i].dur)
				cb.SetLimit(runTable[i].limit)
				cb.SetBurst(runTable[i].burst)
				//cb = table[i].cb //NewStatusCircuitBreaker(table[i].limit, table[i].burst, target.Select())
				limiterTime = time.Now().UTC()
			}
		default:
		}
		select {
		case <-quit:
			//c := fmt.Sprintf("received quit with ping count: %v\n", count)
			status <- runtime.NewStatusOK().SetContent(fmt.Sprintf("quit -> elapsed time: %v", time.Since(start)), false)
			return
		default:
		}
	}
}
