package resiliency

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

var okCircuitBreaker = NewStatusCircuitBreaker(200, 100, func(s *runtime.Status) bool { return s.OK() })

var okPing = func(ctx context.Context) *runtime.Status { return runtime.NewStatusOK() }

func Example_exponentialDecay() {
	//var ms float64 = 1000.0
	rate := .4
	y := 5000.0
	for i := 0; i < 10; i++ {
		if i > 0 {
			y = y * (1 - rate)
		}
		fmt.Printf("test: expDecay() -> %v = %v\n", i, y)
	}
	//y = y * (1 - rate)
	//fmt.Printf("test: expDecay() -> %v\n", y)

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

	go runTicks(okPing, okCircuitBreaker, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runTicks() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}

func Example_runPing() {
	quit := make(chan struct{})
	status := make(chan *runtime.Status, 100)

	go runPing(okPing, okCircuitBreaker, quit, status)
	time.Sleep(time.Millisecond * 500)
	quit <- struct{}{}
	s := <-status

	fmt.Printf("test: runPing() -> [status:%v] %v\n", s, s.ContentString())

	//Output:

}
