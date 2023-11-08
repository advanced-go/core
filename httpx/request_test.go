package httpx

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/log"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

func Example_NewRequest_Nil() {
	newReq, status := NewRequest(nil, "get", "https://www/google.com/search?q=golang", "variant:location")
	fn := log.AccessFromAny(newReq)
	fmt.Printf("test: NewRequest(nil) -> [status:%v] [req-len:%v] [ctx-len:%v] [access:%v] [var:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())), fn != nil, newReq.Header.Get(runtime.ContentLocation))

	//Output:
	//test: NewRequest(nil) -> [status:OK] [req-len:36] [ctx-len:36] [access:true] [var:variant:location]

}

func Example_NewRequest_Context() {
	ctx := context.Background()
	newReq, status := NewRequest(ctx, "get", "https://www/google.com/search?q=golang", "variant:location-2")
	fn := log.AccessFromAny(newReq)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v] [access:%v] [var:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())), fn != nil, newReq.Header.Get(runtime.ContentLocation))

	ctx = runtime.NewRequestIdContext(context.Background(), "123456")
	newReq, status = NewRequest(ctx, "get", "https://www/google.com/search?q=golang", "variant:location")
	fn = log.AccessFromAny(newReq)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v] [access:%v] [var:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())), fn != nil, newReq.Header.Get(runtime.ContentLocation))

	//Output:
	//test: NewRequest(ctx) -> [status:OK] [req-len:36] [ctx-len:36] [access:true] [var:variant:location-2]
	//test: NewRequest(ctx) -> [status:OK] [req-len:6] [ctx-len:6] [access:true] [var:variant:location]

}

func Example_NewRequest_Request() {
	//ctx := context.Background()

	req, _ := http.NewRequest("", "https://www/google.com/search?q=golang", nil)
	newReq, status := NewRequest(req, "get", "https://www/google.com/search?q=golang", "variant:location-2")
	fn := log.AccessFromAny(newReq)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v] [access:%v] [var:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())), fn != nil, newReq.Header.Get(runtime.ContentLocation))

	req, _ = http.NewRequest("", "https://www/google.com/search?q=golang", nil)
	req.Header.Add(runtime.XRequestId, "1234-5678")
	newReq, status = NewRequest(req, "get", "https://www/google.com/search?q=golang", "variant:location-3")
	fn = log.AccessFromAny(newReq)
	fmt.Printf("test: NewRequest(ctx) -> [status:%v] [req-len:%v] [ctx-len:%v] [access:%v] [var:%v]\n", status, len(runtime.RequestId(newReq)), len(runtime.RequestId(newReq.Context())), fn != nil, newReq.Header.Get(runtime.ContentLocation))

	//Output:
	//test: NewRequest(ctx) -> [status:OK] [req-len:36] [ctx-len:36] [access:true] [var:variant:location-2]
	//test: NewRequest(ctx) -> [status:OK] [req-len:9] [ctx-len:9] [access:true] [var:variant:location-3]

}
