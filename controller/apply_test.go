package controller

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/exchange"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func ExampleApply_SameContext() {
	uri := "https://www.google.com/search?q=golang"
	h := make(http.Header)
	ctx := context.Background()
	status := runtime.StatusOK()
	var newCtx context.Context

	defer Apply(ctx, &newCtx, access.NewRequest(h, http.MethodGet, uri), nil, "google-search", 0, access.StatusCode(&status))()
	fmt.Printf("test: Apply(\"0ms\") -> [ctx==newCtx:%v]\n", ctx == newCtx)

	ctx1, cancel := context.WithTimeout(ctx, time.Millisecond*333)
	defer cancel()
	defer Apply(ctx1, &newCtx, access.NewRequest(h, http.MethodGet, uri), nil, "google-search", time.Millisecond*100, access.StatusCode(&status))()
	fmt.Printf("test: Apply(\"100ms\") -> [ctx==newCtx:%v]\n", ctx1 == newCtx)

	//Output:
	//test: Apply("0ms") -> [ctx==newCtx:true]
	//test: Apply("100ms") -> [ctx==newCtx:true]

}

func ExampleApply_NewContext() {
	uri := "https://www.google.com/search?q=golang"
	h := make(http.Header)
	status := runtime.StatusOK()
	var newCtx context.Context

	ctx := context.Background()
	defer Apply(ctx, &newCtx, access.NewRequest(h, http.MethodGet, uri), nil, "google-search", time.Millisecond*100, access.StatusCode(&status))()
	fmt.Printf("test: Apply(\"0ms\") -> [ctx==newCtx:%v]\n", ctx == newCtx)

	//Output:
	//test: Apply("0ms") -> [ctx==newCtx:false]

}

func ExampleApply_Timeout_1000ms() {
	uri := "https://www.google.com/search?q=golang"
	h := make(http.Header)
	var newCtx context.Context
	var resp *http.Response
	var status *runtime.Status

	defer Apply(nil, &newCtx, access.NewRequest(h, http.MethodGet, uri), &resp, "google-search", time.Millisecond*1000, access.StatusCode(&status))()
	resp, status = exchange.Get(newCtx, uri, h)
	if status.OK() {
		buf, _ := io2.ReadAll(resp.Body, h)
		resp.ContentLength = int64(len(buf))
	}
	fmt.Printf("test: exchange.Get(\"1000ms\") -> [status:%v] [status-code:%v] [content-type:%v]\n", status, resp.StatusCode, resp.Header.Get("Content-Type"))

	//Output:
	//test: exchange.Get("1000ms") -> [status:OK] [status-code:200] [content-type:text/html; charset=ISO-8859-1]

}

func ExampleApply_Timeout_10ms() {
	uri := "https://www.google.com/search?q=golang"
	h := make(http.Header)
	var newCtx context.Context
	var resp *http.Response
	var status *runtime.Status

	defer Apply(nil, &newCtx, access.NewRequest(h, http.MethodGet, uri), &resp, "google-search", time.Millisecond*10, access.StatusCode(&status))()
	resp, status = exchange.Get(newCtx, uri, h)
	fmt.Printf("test: exchange.Get(\"10ms\") -> [status:%v] [status-code:%v] [content-type:%v]\n", status, resp.StatusCode, resp.Header.Get("Content-Type"))

	//Output:
	//test: exchange.Get("10ms") -> [status:Deadline Exceeded [Get "https://www.google.com/search?q=golang": context deadline exceeded]] [status-code:4] [content-type:]

}

func ExampleCreateResponse() {
	var r *http.Response

	resp := createResponse(nil, http.StatusOK)
	fmt.Printf("test: createResponse(nil) -> [status-code:%v] [status:%v]\n", resp.StatusCode, resp.Status)

	resp = createResponse(&r, http.StatusGatewayTimeout)
	fmt.Printf("test: createResponse(nil) -> [status-code:%v] [status:%v]\n", resp.StatusCode, resp.Status)

	r = new(http.Response)
	r.StatusCode = http.StatusTeapot
	r.Status = "I'm a Teapot"
	resp = createResponse(&r, http.StatusGatewayTimeout)
	fmt.Printf("test: createResponse(nil) -> [status-code:%v] [status:%v]\n", resp.StatusCode, resp.Status)

	//Output:
	//test: createResponse(nil) -> [status-code:200] [status:OK]
	//test: createResponse(nil) -> [status-code:504] [status:Timeout]
	//test: createResponse(nil) -> [status-code:418] [status:I'm a Teapot]

}
