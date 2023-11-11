package http2

import (
	"context"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"reflect"
)

func Example_Do() {
	ctx := context.Background()
	do(ctx, "GET", "https://www.google.com/search?q=golang", nil)

	req, _ := http.NewRequest("put", "https://twitter.com", nil)
	req.Header.Add("request-id", "1234-5678")
	do(req, "GET", "https://www.google.com/search?q=golang", nil)

	//Output:
}

func do(ctx any, method, uri string, body any) (any, *runtime.Status) {
	fmt.Printf("test: do() -> [type:%v] [method:%v] [uri:%v]\n", reflect.TypeOf(ctx), method, uri)
	return nil, nil
}
