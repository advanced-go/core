package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

var okPing = func(ctx context.Context) *runtime.Status { return runtime.NewStatusOK() }

func createTable() []runArgs {
	var table []runArgs
	for _, arg := range runTable {
		table = append(table, arg)
	}
	return table
}

func Example_StatusAgent_Error() {
	cb, _ := NewStatusCircuitBreaker(Threshold{Limit: 100, Burst: 50, Duration: 0, Select: okSelect})

	_, err := NewStatusAgent(-1, nil, cb)
	fmt.Printf("test: NewStatusAgent() -> %v\n", err)

	_, err = NewStatusAgent(-1, okPing, nil)
	fmt.Printf("test: NewStatusAgent() -> %v\n", err)

	//Output:
	//test: NewStatusAgent() -> error: ping function is nil
	//test: NewStatusAgent() -> error: circuit breaker is nil

}

func Example_runTest() {
	idiomaticGo := true
	quit := make(chan struct{}, 1)
	status := make(chan *runtime.Status, 100)
	cb, _ := NewStatusCircuitBreaker(Threshold{Limit: 100, Burst: 100, Duration: 0, Select: func(s *runtime.Status) bool { return s.OK() }})

	go run(createTable(), func(ctx context.Context) *runtime.Status { return runtime.NewStatusOK() }, 0, cb, quit, status)
	if idiomaticGo {
		done := make(chan struct{})
		go func(chan struct{}, chan *runtime.Status) {
			for {
				select {
				case st := <-status:
					if st.IsContent() {
						fmt.Printf("test: runTest() -> %v", st.ContentString())
					}
					if st.OK() {
						done <- struct{}{}
						return
					}
				default:
				}
			}
		}(done, status)
		<-done
		close(done)
	} else {
		// idiomatic Java
		time.Sleep(time.Minute * 1)
		quit <- struct{}{}
		s := <-status
		if s != nil && s.IsContent() {
			fmt.Printf("test: runTest() -> %v\n", s.ContentString())
		}
		//for s := range status {
		//		}
		//	}
	}
	close(quit)
	close(status)

	//Output:

}

/*
func _Example_runChannels() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)

	go runChannels(okPing, nil, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runChannels() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func _Example_runTicks() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)
	_,cb := NewStatusCircuitBreaker(200, 200, func(s *runtime.Status) bool { return s.OK() })

	go runTicks(okPing, cb, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runTicks() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func runStatus(_ PingFn, _[]runArgs, _ StatusCircuitBreaker, quit chan struct{}, status chan *runtime.Status) {
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


*/
