package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

var okCircuitBreaker = NewStatusCircuitBreaker(200, 200, func(s *runtime.Status) bool { return s.OK() })

var okPing = func(ctx context.Context) *runtime.Status { return runtime.NewStatusOK() }

func _Example_NewStatusAgent() {
	a := NewStatusAgent(0, okPing, okCircuitBreaker)
	fmt.Printf("test: NewStatusAgent() -> %v\n", a)

	//Output:
}

func _Example_runChannels() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)

	go runChannels(okPing, okCircuitBreaker, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runChannels() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func _Example_runTicks() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)
	cb := NewStatusCircuitBreaker(200, 200, func(s *runtime.Status) bool { return s.OK() })

	go runTicks(okPing, cb, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runTicks() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func Example_runTest() {
	quit := make(chan struct{}, 1)
	status := make(chan *runtime.Status, 100)
	cb := NewStatusCircuitBreaker(100, 10, func(s *runtime.Status) bool { return s.OK() })

	go runTest(createTable(), okPing, 0, cb, quit, status)
	time.Sleep(time.Minute * 2)
	quit <- struct{}{}
	s := <-status
	close(quit)
	close(status)

	fmt.Printf("test: runPing() -> %v\n", s.ContentString())

	//Output:

}

func _Example_runSlice() {
	quit := make(chan struct{}, 1)
	status := make(chan *runtime.Status, 100)
	cb := NewStatusCircuitBreaker(100, 10, func(s *runtime.Status) bool { return s.OK() })

	args := []runArgs{{limit: rate.Limit(90), burst: 100, dur: time.Millisecond * 1}}

	go runTest(args, okPing, 0, cb, quit, status)
	time.Sleep(time.Minute * 1)
	quit <- struct{}{}
	s := <-status
	close(quit)
	close(status)

	fmt.Printf("test: runPing() -> %v\n", s.ContentString())

	//Output:

}
