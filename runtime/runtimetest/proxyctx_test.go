package runtimetest

import (
	"context"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func testProxyContext(ctx context.Context) bool {
	if _, ok := any(ctx).(*proxyContext); ok {
		return true
	}
	return false
}

func proxyGet(uri string) (*http.Response, error) {
	fmt.Printf("test: testDo() -> \n")
	return nil, errors.New("test error")
}

func proxyDo(req *http.Request) (*http.Response, error) {
	fmt.Printf("test: testDo() -> \n")
	return nil, errors.New("test error")
}

func ExampleProxyContext() {
	k1 := "1"
	k2 := "2"
	k3 := "3"
	v1 := "value 1"
	v2 := "value 2"
	v3 := "value 3"

	ctx := runtime.NewProxyContext(nil, proxyDo)

	fmt.Printf("test: isProxyContext(ctx) -> %v\n", testProxyContext(ctx))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx.Value(k1), ctx.Value(k2), ctx.Value(k3))

	ctx1 := runtime.ContextWithValue(ctx, k1, v1)
	fmt.Printf("test: isProxyContext(ctx1) -> %v\n", testProxyContext(ctx1))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx1.Value(k1), ctx1.Value(k2), ctx1.Value(k3))

	ctx2 := runtime.ContextWithValue(ctx, k2, v2)
	fmt.Printf("test: isProxyContext(ctx2) -> %v\n", testProxyContext(ctx2))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx2.Value(k1), ctx2.Value(k2), ctx2.Value(k3))

	ctx3 := runtime.ContextWithValue(ctx, k3, v3)
	fmt.Printf("test: isProxyContext(ctx3) -> %v\n", testProxyContext(ctx3))
	fmt.Printf("test: Values() -> [key1:%v] [key2:%v] [key3:%v]\n", ctx3.Value(k1), ctx3.Value(k2), ctx3.Value(k3))

	//Output:
	//test: isProxyContext(ctx) -> true
	//test: Values() -> [key1:<nil>] [key2:<nil>] [key3:<nil>]
	//test: isProxyContext(ctx1) -> true
	//test: Values() -> [key1:value 1] [key2:<nil>] [key3:<nil>]
	//test: isProxyContext(ctx2) -> true
	//test: Values() -> [key1:value 1] [key2:value 2] [key3:<nil>]
	//test: isProxyContext(ctx3) -> true
	//test: Values() -> [key1:value 1] [key2:value 2] [key3:value 3]

}

func ExampleProxyContext_Proxy() {
	ctx0 := runtime.NewProxyContext(nil, proxyGet)
	ctx := runtime.NewProxyContext(ctx0, proxyDo)
	ok1 := testProxyContext(ctx)

	fmt.Printf("test: isProxyContext(ctx) -> %v\n", ok1)
	if proxies, ok := runtime.IsProxyable(ctx); ok {
		for _, p := range proxies {
			if fn, ok2 := p.(func(*http.Request) (*http.Response, error)); ok2 {
				if fn != nil {
					fmt.Printf("test: proxyDo(*http.Request) -> %v\n", true)
				}
			}
			if fn, ok2 := p.(func(string) (*http.Response, error)); ok2 {
				if fn != nil {
					fmt.Printf("test: proxyGet(*http.Request) -> %v\n", true)
				}
			}
		}
	}

	//Output:
	//test: isProxyContext(ctx) -> true
	//test: proxyGet(*http.Request) -> true
	//test: proxyDo(*http.Request) -> true

}
