package resiliency

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

func _Example_CircuitBreaker() {
	var limit rate.Limit = 0.0
	var burst int = 0
	var fn StatusSelect

	src := NewStatusCircuitBreaker(100, 50, func(status *runtime.Status) bool { return status.OK() })
	if cfg, ok := any(src).(*circuitConfig); ok {
		limit = cfg.limiter.Limit()
		burst = cfg.limiter.Burst()
		fn = cfg.fn
	}
	fmt.Printf("test: CircuitBreaker() -> [limit:%v] [burst:%v] [select:%v]\n", limit, burst, fn)

	/*
		src = CloneStatusCircuitBreaker(src, 45, 15)
		if cfg, ok := any(src).(*circuitConfig); ok {
			limit = cfg.limiter.Limit()
			burst = cfg.limiter.Burst()
			fn2 = cfg.selectFn
		}
		fmt.Printf("test: CloneCircuitBreaker() -> [limit:%v] [burst:%v] [select:%v]\n", limit, burst, fn2)


	*/
	//Output:
	//test: CircuitBreaker() -> [limit:100] [burst:50] [select:0x10202c0]
	//test: CloneCircuitBreaker() -> [limit:45] [burst:15] [select:0x10202c0]

}

func _Example_RateLimiter() {
	rl := rate.NewLimiter(0.5, 1)

	allow := rl.Allow()
	fmt.Printf("test: Allow() -> %v\n", allow)

	//Output:
	//test: Allow() -> true

}

func Example_CircuitTest() {
	count := 1000
	ms := time.Duration(999)

	//testBreaker(1, 1, func(status *runtime.Status) bool { return true }, time.Millisecond*ms, count)
	//testBreaker2(1, 1, time.Millisecond*ms, count)

	// 99ms should work
	ms = time.Duration(95)
	//testBreaker(10, 10, func(status *runtime.Status) bool { return true }, time.Millisecond*ms, count)
	testBreaker2(10, 10, time.Millisecond*ms, count)

	// 19 ms should work
	//ms = time.Duration(15)
	//testBreaker(50, 50, func(status *runtime.Status) bool { return true }, time.Millisecond*ms, count)
	//testBreaker2(50, 50, time.Millisecond*ms, count)

	//Output:
}

func testBreaker(limit rate.Limit, burst int, fn StatusSelect, d time.Duration, count int) {
	start := time.Now().UTC()
	cb := NewStatusCircuitBreaker(limit, burst, fn)
	s := runtime.NewStatusOK()
	for i := 0; i < count; i++ {
		time.Sleep(d)
		if !cb.Allow(s) {
			fmt.Printf("test: testBreaker()  ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "broken", limit, d, i, time.Since(start))
			return
		}
	}
	fmt.Printf("test: testBreaker()  ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "OK", limit, d, count, time.Since(start))
}

func testBreaker2(limit rate.Limit, burst int, d time.Duration, count int) {
	start := time.Now().UTC()
	limiter := rate.NewLimiter(limit, burst)
	for i := 0; i < count; i++ {
		time.Sleep(d)
		if !limiter.Allow() {
			fmt.Printf("test: testBreaker2() ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "broken", limit, d, i, time.Since(start))
			return
		}
	}
	fmt.Printf("test: testBreaker2() ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "OK", limit, d, count, time.Since(start))
}
