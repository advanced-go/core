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

	status := NewStatusOK().SetRequestId("987-654")
	id = RequestId(status)
	fmt.Printf("test: RequestId() -> %v\n", id)

	//Output:
	//test: RequestId() -> 123-456
	//test: RequestId() -> 123-456-abc
	//test: RequestId() -> 123-456-789
	//test: GetOrCreateRequestId() -> [valid:true]
	//test: RequestId() -> 987-654

}

func Example_ContentLocation() {
	ctx := NewContentLocationContext(context.Background(), nil)
	uri, ok := ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h := make(http.Header)
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h.Add(ContentLocation, "https://www.google.com/search?q=golang")
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h = make(http.Header)
	h.Add(ContentLocation, "file://[cwd]/runtimetest/test.txt")
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:true] [uri:file://[cwd]/runtimetest/test.txt]

}

func Example_FileUrl_File() {
	ctx := NewFileUrlContext(context.Background(), "")
	uri, ok := FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url := "https://www.google.com/search?q=golang"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url = "file://[cwd]/runtimetest/test.txt"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:true] [uri:file://[cwd]/runtimetest/test.txt]

}

func Example_FileUrl_Urn() {
	ctx := NewFileUrlContext(context.Background(), "")
	uri, ok := FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url := "https://www.google.com/search?q=golang"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url = "urn:status:ok"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:true] [uri:urn:status:ok]
	
}
