package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

var okCircuitBreaker = NewStatusCircuitBreaker(200, 100, func(s *runtime.Status) bool { return s.OK() })

var okPing = func(ctx context.Context) *runtime.Status { return runtime.NewStatusOK() }

func Example_runChannels() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)

	go runChannels(okPing, okCircuitBreaker, quit, status)

	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	time.Sleep(time.Millisecond * 500)
	s := <-status

	fmt.Printf("test: runChannels() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func Example_runTicks() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)

	go runTicks(okPing, okCircuitBreaker, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	time.Sleep(time.Millisecond * 500)
	s := <-status

	fmt.Printf("test: runTicks() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}
