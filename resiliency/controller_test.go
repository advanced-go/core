package resiliency

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

var handler = func(ctx any, r *http.Request, body any) (any, runtime.Status) {
	return nil, runtime.StatusOK()
}

func Example_Controller() {
	c := NewController("test", Threshold{time.Millisecond * 500}, nil)
	fmt.Printf("test: NewController() -> [err:%v] %v\n", nil, c)

	c = NewController("test", Threshold{time.Millisecond * 500}, handler)
	fmt.Printf("test: NewController() -> [err:%v] %v\n", nil, c)

	//Output:
	//test: NewController() -> [err:error: handler is nil] <nil>
	//test: NewController() -> [err:<nil>] &{test {500000000} 0xc7db80 <nil>}

}
