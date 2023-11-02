package resiliency

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"golang.org/x/time/rate"
)

func _Example_CircuitBreaker() {
	var limit rate.Limit = 0.0
	var burst int = 0
	var fn Select
	var fn2 Select

	src := NewStatusCircuitBreaker(100, 50, func(status *runtime.Status) bool { return status.OK() })
	if cfg, ok := any(src).(*circuitConfig); ok {
		limit = cfg.limiter.Limit()
		burst = cfg.limiter.Burst()
		fn = cfg.selectFn
	}
	fmt.Printf("test: CircuitBreaker() -> [limit:%v] [burst:%v] [select:%v]\n", limit, burst, fn)

	src = CloneStatusCircuitBreaker(src, 45, 15)
	if cfg, ok := any(src).(*circuitConfig); ok {
		limit = cfg.limiter.Limit()
		burst = cfg.limiter.Burst()
		fn2 = cfg.selectFn
	}
	fmt.Printf("test: CloneCircuitBreaker() -> [limit:%v] [burst:%v] [select:%v]\n", limit, burst, fn2)

	//Output:
	//test: CircuitBreaker() -> [limit:100] [burst:50] [select:0x10202c0]
	//test: CloneCircuitBreaker() -> [limit:45] [burst:15] [select:0x10202c0]

}
