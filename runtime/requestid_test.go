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

func ExampleContextWithHeader() {
	ctx := NewRequestIdContext(context.Background(), "123-456-abc")
	fmt.Printf("test: NewRequestIdContext(ctx,id) -> %v\n", RequestIdFromContext(ctx))

	ctx = NewRequestIdContextFromHeader(nil)
	fmt.Printf("test: NewRequestContext(nil) -> %v\n", RequestIdFromContext(ctx) != "")

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	ctx = NewRequestIdContextFromHeader(req.Header)
	fmt.Printf("test: NewRequestContext(req) -> %v\n", RequestIdFromContext(ctx) != "")

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Add(XRequestId, "x-request-id-value")
	ctx = NewRequestIdContextFromHeader(req.Header)
	fmt.Printf("test: NewRequestContext(req) -> %v\n", RequestIdFromContext(ctx))

	//Output:
	//test: NewRequestIdContext(ctx,id) -> 123-456-abc
	//test: NewRequestContext(nil) -> false
	//test: NewRequestContext(req) -> true
	//test: NewRequestContext(req) -> x-request-id-value

}

func ExampleRequestId() {
	id := RequestId("123-456-string")
	fmt.Printf("test: RequestId(string) -> %v\n", id)

	ctx := NewRequestIdContext(context.Background(), "123-456-context")
	id = RequestId(ctx)
	fmt.Printf("test: RequestId(ctx) -> %v\n", id)

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-request")
	id = RequestId(req)
	fmt.Printf("test: RequestId(request) -> %v\n", id)

	h := make(http.Header)
	h.Set(XRequestId, "123-456-header")
	id = RequestId(h)
	fmt.Printf("test: RequestId(header) -> %v\n", id)

	//Output:
	//test: RequestId(string) -> 123-456-string
	//test: RequestId(ctx) -> 123-456-context
	//test: RequestId(request) -> 123-456-request
	//test: RequestId(header) -> 123-456-header

}

func ExampleRequestId_New() {
	id := RequestId(nil)
	fmt.Printf("test: RequestId(nil) -> [empty:%v]\n", len(id) == 0)

	id = RequestId(context.Background())
	fmt.Printf("test: RequestId(ctx) -> [empty:%v]\n", len(id) == 0)

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	id = RequestId(req)
	fmt.Printf("test: RequestId(request) -> [empty:%v]\n", len(id) == 0)

	h := make(http.Header)
	id = RequestId(h)
	fmt.Printf("test: RequestId(header) -> [empty:%v]\n", len(id) == 0)

	//Output:
	//test: RequestId(nil) -> [empty:false]
	//test: RequestId(ctx) -> [empty:false]
	//test: RequestId(request) -> [empty:false]
	//test: RequestId(header) -> [empty:false]

}

func ExampleAddRequestId() {
	h := AddRequestId(nil)
	fmt.Printf("test: AddRequestId(nil) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	head := make(http.Header)
	h = AddRequestId(head)
	fmt.Printf("test: AddRequestId(head) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	head = make(http.Header)
	head.Add(XRequestId, "123-45-head")
	h = AddRequestId(head)
	fmt.Printf("test: AddRequestId(head) -> %v\n", h.Get(XRequestId))

	req, _ := http.NewRequest("", "https.www.google.com", nil)
	h = AddRequestId(req)
	fmt.Printf("test: RequestId(request) -> [empty:%v]\n", len(h.Get(XRequestId)) == 0)

	req, _ = http.NewRequest("", "https.www.google.com", nil)
	req.Header.Set(XRequestId, "123-456-request")
	h = AddRequestId(req)
	fmt.Printf("test: RequestId(request) -> %v\n", h.Get(XRequestId))

	//Output:
	//test: AddRequestId(nil) -> [empty:false]
	//test: AddRequestId(head) -> [empty:false]
	//test: AddRequestId(head) -> 123-45-head
	//test: RequestId(request) -> [empty:false]
	//test: RequestId(request) -> 123-456-request
	
}
