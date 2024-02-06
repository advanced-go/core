package controller

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/exchange"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func _ExampleCreateRequest() {
	req, _ := http.NewRequest(http.MethodPut, "https://www.google.com/search?q=golang", nil)

	fmt.Printf("test: CreateRequest() -> [method:%v] [uri:%v]\n", req.Method, req.URL.String())

	//Output:
}

func ExampleApply_Timeout_1000ms() {
	uri := "https://www.google.com/search?q=golang"
	h := make(http.Header)
	var newCtx context.Context
	var resp *http.Response
	var status *runtime.Status

	defer Apply(nil, &newCtx, http.MethodGet, uri, "google-search", h, time.Millisecond*1000, access.StatusCode(&status))()
	resp, status = exchange.Get(newCtx, uri, h)
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

	defer Apply(nil, &newCtx, http.MethodGet, uri, "google-search", h, time.Millisecond*10, access.StatusCode(&status))()
	resp, status = exchange.Get(newCtx, uri, h)
	fmt.Printf("test: exchange.Get(\"10ms\") -> [status:%v] [status-code:%v] [content-type:%v]\n", status, resp.StatusCode, resp.Header.Get("Content-Type"))

	//Output:
	//test: exchange.Get("10ms") -> [status:Deadline Exceeded [Get "https://www.google.com/search?q=golang": context deadline exceeded]] [status-code:4] [content-type:]

}
