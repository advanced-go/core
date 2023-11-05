package log

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var testLogger = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string) {
}

func ExampleContextWithAccessLoggerExisting() {
	ctx := ContextWithAccessLogger(context.Background(), testLogger)
	fmt.Printf("test: ContextWithAccessLogger(context.Background(),id) -> %v [newContext:%v]\n", ContextAccessLogger(ctx), ctx != context.Background())

	ctxNew := ContextWithAccessLogger(ctx, testLogger)
	fmt.Printf("test: ContextWithAccessLogger(ctx,id) -> %v [newContext:%v]\n", ContextAccessLogger(ctx), ctxNew != ctx)

	//Output:
	//test: ContextWithAccessLogger(context.Background(),id) -> 123-456-abc [newContext:true]
	//test: ContextWithAccessLogger(ctx,id) -> 123-456-abc [newContext:false]

}

/*
func ExampleContextWithAccessLogger() {
	ctx := ContextWithAccessLogger(context.Background(), testLogger)
	fmt.Printf("test: ContextWithAccessLogger(ctx,id) -> %v\n", ContextAccessLogger(ctx))

	ctx = ContextWithAccessLogger(nil)
	fmt.Printf("test: ContextWithAccessLogger(nil) -> %v\n", ContextAccessLogger(ctx) != "")

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	ctx = ContextWithRequest(req)
	fmt.Printf("test: ContextWithRequest(req) -> %v\n", ContextAccessLogger(ctx) != "")

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Add(XRequestId, "x-request-id-value")
	ctx = ContextWithRequest(req)
	fmt.Printf("test: ContextWithRequest(req) -> %v\n", ContextAccessLogger(ctx))

	//Output:
	//test: ContextWithAccessLogger(ctx,id) -> 123-456-abc
	//test: ContextWithRequest(nil) -> false
	//test: ContextWithRequest(req) -> true
	//test: ContextWithRequest(req) -> x-request-id-value

}

func Example_RequestId() {
	id := RequestId("123-456")
	fmt.Printf("test: RequestId() -> %v\n", id)

	ctx := ContextWithAccessLogger(context.Background(), "123-456-abc")
	id = RequestId(ctx)
	fmt.Printf("test: RequestId() -> %v\n", id)

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-789")
	id = RequestId(req)
	fmt.Printf("test: RequestId() -> %v\n", id)

	status := NewStatusOK().SetRequestId("987-654")
	id = RequestId(status)
	fmt.Printf("test: RequestId() -> %v\n", id)

	//Output:
	//test: RequestId() -> 123-456
	//test: RequestId() -> 123-456-abc
	//test: RequestId() -> 123-456-789
	//test: RequestId() -> 987-654

}


*/
