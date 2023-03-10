package runtime

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

func isExchangeContext(ctx context.Context) bool {
	if _, ok := any(ctx).(withValue); ok {
		return true
	}
	return false
}

func testDo(req *http.Request) (*http.Response, error) {
	fmt.Printf("test: testDo() -> \n")
	return nil, errors.New("test error")
}

func ExampleExchangeContext() {
	k1 := "1"
	k2 := "2"
	k3 := "3"
	v1 := "value 1"
	v2 := "value 2"
	v3 := "value 3"

	ctx := ContextWithHttpExchange(nil, testDo)

	fmt.Printf("test: isExchangeContext(ctx) -> %v\n", isExchangeContext(ctx))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx.Value(k1), ctx.Value(k2), ctx.Value(k3))

	ctx1 := ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: isExchangeContext(ctx1) -> %v\n", isExchangeContext(ctx1))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx1.Value(k1), ctx1.Value(k2), ctx1.Value(k3))

	ctx2 := ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: isExchangeContext(ctx2) -> %v\n", isExchangeContext(ctx2))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx2.Value(k1), ctx2.Value(k2), ctx2.Value(k3))

	ctx3 := ContextWithValue(ctx, k3, v3)
	fmt.Printf("test: isExchangeContext(ctx3) -> %v\n", isExchangeContext(ctx3))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx3.Value(k1), ctx3.Value(k2), ctx3.Value(k3))

	//Output:
	//test: isExchangeContext(ctx) -> true
	//test: Values() -> [key1:<nil>] [key2:<nil>] [key3:<nil>]
	//test: isExchangeContext(ctx1) -> true
	//test: Values() -> [key1:value 1] [key2:<nil>] [key3:<nil>]
	//test: isExchangeContext(ctx2) -> true
	//test: Values() -> [key1:value 1] [key2:value 2] [key3:<nil>]
	//test: isExchangeContext(ctx3) -> true
	//test: Values() -> [key1:value 1] [key2:value 2] [key3:value 3]

}
