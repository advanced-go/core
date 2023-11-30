package http2

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func Example_NewRequest_Nil() {
	newReq, status := NewRequest(nil, "get", "https://www/google.com/search?q=golang", nil)
	fmt.Printf("test: NewRequest(nil) -> [status:%v] [req-len:%v] [ctx-len:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())))

	//Output:
	//test: NewRequest(nil) -> [status:OK] [req-len:36] [ctx-len:36]

}

func Example_NewRequest_Context() {
	ctx := context.Background()
	newReq, status := NewRequest(ctx, "get", "https://www/google.com/search?q=golang", nil)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())))

	ctx = runtime.NewRequestIdContext(context.Background(), "123456")
	newReq, status = NewRequest(ctx, "get", "https://www/google.com/search?q=golang", nil)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())))

	//Output:
	//test: NewRequest(ctx) -> [status:OK] [req-len:36] [ctx-len:36]
	//test: NewRequest(ctx) -> [status:OK] [req-len:6] [ctx-len:6]

}

func Example_NewRequest_Request() {
	//ctx := context.Background()

	req, _ := http.NewRequest("", "https://www/google.com/search?q=golang", nil)
	newReq, status := NewRequest(req, "get", "https://www/google.com/search?q=golang", nil)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())))

	req, _ = http.NewRequest("", "https://www/google.com/search?q=golang", nil)
	req.Header.Add(runtime.XRequestId, "1234-5678")
	newReq, status = NewRequest(req, "get", "https://www/google.com/search?q=golang", nil)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())))

	//Output:
	//test: NewRequest(ctx) -> [status:OK] [req-len:36] [ctx-len:36]
	//test: NewRequest(ctx) -> [status:OK] [req-len:9] [ctx-len:9]

}

func Example_Clone() {
	req, _ := http.NewRequest("get", "http://localhost:8080/search?q=golang", nil)
	clone := req.Clone(context.Background())

	fmt.Printf("test: Clone() ->  [orig:%v] [clone:%v] [orig:%v] [clone:%v]\n", req.Host, clone.Host, req.URL.String(), clone.URL.String())
	fmt.Printf("test: Clone() ->  [origUrl:%v] [cloneUrl:%v]\n", req.URL.Host, clone.URL.Host)

	//Output:
	//test: Clone() ->  [orig:localhost:8080] [clone:localhost:8080] [orig:http://localhost:8080/search?q=golang] [clone:http://localhost:8080/search?q=golang]
	//test: Clone() ->  [origUrl:localhost:8080] [cloneUrl:localhost:8080]

}

/*
func Example_NewRequest() {
	req, status := NewRequest(nil, "PUT", "https://somedomain.com/invalid-uri-or-type", "")

	fmt.Printf("test: NewRequest() [status:%v] %v\n", status, req)

	//Output:

}


*/
