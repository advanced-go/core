package log

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

var testLogger = func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, threshold int, statusFlags string) {
	fmt.Printf("test: testLogger() -> %v", "{ access logging attributes }")

}

func Example_ContextWithAccessLoggerExisting() {
	SetAccess(testLogger)
	ctx := NewAccessContext(context.Background())
	fmt.Printf("test: NewAccessContext(context.Background(),id) -> [access:%v] [newContext:%v]\n", AccessFromContext(ctx) != nil, ctx != context.Background())

	ctxNew := NewAccessContext(ctx)
	fmt.Printf("test: NewAccessContext(ctx,id) -> [access:%v] [newContext:%v]\n", AccessFromContext(ctx) != nil, ctxNew != ctx)

	//Output:
	//test: NewAccessContext(context.Background(),id) -> [access:true] [newContext:true]
	//test: NewAccessContext(ctx,id) -> [access:true] [newContext:false]

}

func Example_AccessLogger() {
	start := time.Now().UTC()
	SetAccess(testLogger)
	ctx := NewAccessContext(context.Background())
	if fn := AccessFromContext(ctx); fn != nil {
		fn("egress", start, time.Since(start), nil, nil, -1, "flags")
	}

	//Output:
	//test: testLogger() -> { access logging attributes }

}

/*
func ExampleNewAccessContext() {
	ctx := NewAccessContext(context.Background(), testLogger)
	fmt.Printf("test: NewAccessContext(ctx,id) -> %v\n", AccessFnFromContext(ctx))

	ctx = NewAccessContext(nil)
	fmt.Printf("test: NewAccessContext(nil) -> %v\n", AccessFnFromContext(ctx) != "")

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	ctx = ContextWithRequest(req)
	fmt.Printf("test: ContextWithRequest(req) -> %v\n", AccessFnFromContext(ctx) != "")

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Add(XRequestId, "x-request-id-value")
	ctx = ContextWithRequest(req)
	fmt.Printf("test: ContextWithRequest(req) -> %v\n", AccessFnFromContext(ctx))

	//Output:
	//test: NewAccessContext(ctx,id) -> 123-456-abc
	//test: ContextWithRequest(nil) -> false
	//test: ContextWithRequest(req) -> true
	//test: ContextWithRequest(req) -> x-request-id-value

}

func Example_RequestId() {
	id := RequestId("123-456")
	fmt.Printf("test: RequestId() -> %v\n", id)

	ctx := NewAccessContext(context.Background(), "123-456-abc")
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
