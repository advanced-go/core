package resiliency

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
	"time"
)

var okSelect = func(status *runtime.Status) bool { return status.OK() }

func Example_CircuitBreaker_Error() {
	_, err := NewStatusCircuitBreaker(Threshold{Limit: 0, Burst: 50, Timeout: 0, Select: okSelect})
	fmt.Printf("test: NewStatusCircuitBreaker() -> %v\n", err)

	_, err = NewStatusCircuitBreaker(Threshold{Limit: 100, Burst: 0, Timeout: 0, Select: okSelect})
	fmt.Printf("test: NewStatusCircuitBreaker() -> %v\n", err)

	_, err = NewStatusCircuitBreaker(Threshold{Limit: -1, Burst: 50, Timeout: 0, Select: okSelect})
	fmt.Printf("test: NewStatusCircuitBreaker() -> %v\n", err)

	_, err = NewStatusCircuitBreaker(Threshold{Limit: 101, Burst: 50, Timeout: 0, Select: nil})
	fmt.Printf("test: NewStatusCircuitBreaker() -> %v\n", err)

	_, err = NewStatusCircuitBreaker(Threshold{Limit: 100, Burst: 50, Timeout: 0, Select: nil})
	fmt.Printf("test: NewStatusCircuitBreaker() -> %v\n", err)

	//Output:
	//test: NewStatusCircuitBreaker() -> error: rate limit or burst is invalid limit = 0 burst = 50
	//test: NewStatusCircuitBreaker() -> error: rate limit or burst is invalid limit = 100 burst = 0
	//test: NewStatusCircuitBreaker() -> error: rate limit or burst is invalid limit = -1 burst = 50
	//test: NewStatusCircuitBreaker() -> error: rate limit [101] is greater than the maximum [100]
	//test: NewStatusCircuitBreaker() -> error: status select function in nil

}

func Example_CircuitBreaker_Clone() {
	cb, _ := NewStatusCircuitBreaker(Threshold{Limit: 100, Burst: 50, Timeout: 0, Select: okSelect})
	clone := CloneStatusCircuitBreaker(cb)

	fmt.Printf("test: CloneStatusCircuitBreaker() -> [limit:%v] [burst:%v]\n", clone.Limit(), clone.Burst())

	//Output:
	//test: CloneStatusCircuitBreaker() -> [limit:100] [burst:50]

}

func _Example_CircuitTest() {
	count := 1000
	ms := time.Duration(999)
	limiter := rate.NewLimiter(1, 1)

	testBreaker2(limiter, 1, 1, time.Millisecond*ms, count)

	// 100ms should work  actual 94
	ms = time.Duration(94)
	testBreaker2(limiter, 10, 10, time.Millisecond*ms, count)

	// 40ms should work  actual
	ms = time.Duration(30)
	testBreaker2(limiter, 25, 25, time.Millisecond*ms, count)

	// 20 ms should work actual
	ms = time.Duration(15)
	testBreaker2(limiter, 50, 50, time.Millisecond*ms, count)

	// 13 ms should work actual
	ms = time.Duration(6)
	testBreaker2(limiter, 75, 75, time.Millisecond*ms, count)

	// 1 ms should work actual 0
	ms = time.Duration(0)
	testBreaker2(limiter, 100, 100, time.Millisecond*ms, count+1000)

	//Output:
}

func testBreaker(limit rate.Limit, burst int, fn StatusSelectFn, d time.Duration, count int) {
	start := time.Now().UTC()
	cb, _ := NewStatusCircuitBreaker(Threshold{Limit: limit, Burst: burst, Timeout: 0, Select: fn})
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

func testBreaker2(limiter *rate.Limiter, limit rate.Limit, burst int, d time.Duration, count int) {
	start := time.Now().UTC()
	limiter.SetLimit(limit)
	limiter.SetBurst(burst)
	for i := 0; i < count; i++ {
		if d > 0 {
			time.Sleep(d)
		}
		if !limiter.Allow() {
			fmt.Printf("test: testBreaker2() ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "broken", limit, d, i, time.Since(start))
			return
		}
	}
	fmt.Printf("test: testBreaker2() ->  [circuit:%v] [limit:%v] [duration:%v] [count:%v] [elapsed:%v]\n", "OK", limit, d, count, time.Since(start))
}
