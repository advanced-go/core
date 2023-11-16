package http2

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"reflect"
)

func Example_DoContext() {
	ctx := context.Background()
	do(ctx, "GET", "https://www.google.com/search?q=golang", nil)

	req, _ := http.NewRequest("put", "https://twitter.com", nil)
	req.Header.Add("request-id", "1234-5678")
	do(req, "GET", "https://www.google.com/search?q=golang", nil)

	//Output:
	//test: do() -> [type:context.backgroundCtx] [method:GET] [uri:https://www.google.com/search?q=golang]
	//test: do() -> [type:*http.Request] [method:GET] [uri:https://www.google.com/search?q=golang]

}

func do(ctx any, method, uri string, body any) (any, *runtime.Status) {
	fmt.Printf("test: do() -> [type:%v] [method:%v] [uri:%v]\n", reflect.TypeOf(ctx), method, uri)
	return nil, nil
}

func doHandler(ctx any, r *http.Request, body any) (any, *runtime.Status) {
	return nil, runtime.NewStatusOK()
}

func Example_DoHandlerProxy() {
	ctx := runtime.NewProxyContext(nil, doHandler)

	fn := DoHandlerProxy(ctx)
	fmt.Printf("test: DoHandlerProxy(ctx) -> [proxy:%v]\n", fn != nil)

	req, _ := http.NewRequestWithContext(ctx, "", "https://www.google.com/search", nil)
	fn = DoHandlerProxy(req)
	fmt.Printf("test: DoHandlerProxy(*http.Request) -> [proxy:%v]\n", fn != nil)

	//Output:
	//test: DoHandlerProxy(ctx) -> [proxy:true]
	//test: DoHandlerProxy(*http.Request) -> [proxy:true]

}
