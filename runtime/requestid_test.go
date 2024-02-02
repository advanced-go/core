package runtime

import (
	"context"
	"fmt"
	"net/http"
)

func ExampleContextWithRequestExisting() {
	ctx := NewRequestIdContext(context.Background(), "123-456-abc")
	fmt.Printf("test: NewRequestIdContext(context.Background(),id) -> %v [newContext:%v]\n", RequestIdFromContext(ctx), ctx != context.Background())

	ctxNew := NewRequestIdContext(ctx, "123-456-abc-xyz")
	fmt.Printf("test: NewRequestIdContext(ctx,id) -> %v [newContext:%v]\n", RequestIdFromContext(ctx), ctxNew != ctx)

	//Output:
	//test: NewRequestIdContext(context.Background(),id) -> 123-456-abc [newContext:true]
	//test: NewRequestIdContext(ctx,id) -> 123-456-abc [newContext:false]

}

func ExampleContextWithRequest() {
	ctx := NewRequestIdContext(context.Background(), "123-456-abc")
	fmt.Printf("test: NewRequestIdContext(ctx,id) -> %v\n", RequestIdFromContext(ctx))

	ctx = NewRequestContext(nil)
	fmt.Printf("test: NewRequestContext(nil) -> %v\n", RequestIdFromContext(ctx) != "")

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	ctx = NewRequestContext(req)
	fmt.Printf("test: NewRequestContext(req) -> %v\n", RequestIdFromContext(ctx) != "")

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Add(XRequestId, "x-request-id-value")
	ctx = NewRequestContext(req)
	fmt.Printf("test: NewRequestContext(req) -> %v\n", RequestIdFromContext(ctx))

	//Output:
	//test: NewRequestIdContext(ctx,id) -> 123-456-abc
	//test: NewRequestContext(nil) -> false
	//test: NewRequestContext(req) -> true
	//test: NewRequestContext(req) -> x-request-id-value

}

func Example_RequestId() {
	id := RequestId("123-456")
	fmt.Printf("test: RequestId() -> %v\n", id)

	ctx := NewRequestIdContext(context.Background(), "123-456-abc")
	id = RequestId(ctx)
	fmt.Printf("test: RequestId() -> %v\n", id)

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-789")
	id = RequestId(req)
	fmt.Printf("test: RequestId() -> %v\n", id)

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	id = GetOrCreateRequestId(req)
	if req.Header.Get(XRequestId) == "" {
		req.Header.Set(XRequestId, id)
	}
	id = RequestId(req)
	fmt.Printf("test: GetOrCreateRequestId() -> [valid:%v]\n", len(id) != 0)

	//status := NewStatusOK().SetRequestId("987-654")
	//id = RequestId(status)
	//fmt.Printf("test: RequestId() -> %v\n", id)

	//Output:
	//test: RequestId() -> 123-456
	//test: RequestId() -> 123-456-abc
	//test: RequestId() -> 123-456-789
	//test: GetOrCreateRequestId() -> [valid:true]

}
