package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func ExampleNewResource() {
	uri := "/search?q=golang"
	rsc := NewPrimaryResource("http://localhost:8080", "/health/liveness", 0, httpCall)

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	url := rsc.BuildUri(req.URL)
	fmt.Printf("test: NewResource(\"%v\") [url:%v]\n", uri, url)

	//Output:
	//test: NewResource("/search?q=golang") [url:http://localhost:8080/search?q=golang]

}

func ExampleTimeout() {
	var dIn time.Duration = -1
	uri := "/search?q=golang"
	rsc := NewPrimaryResource("http://localhost:8080", "/health/liveness", dIn, httpCall)

	dOut := rsc.timeout(nil)
	fmt.Printf("test: timeout(nil) -> [timeout:%v]\n", dOut)

	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [timeout:%v]\n", dIn, dOut)

	dIn = time.Millisecond * 100
	rsc = NewPrimaryResource("http://localhost:8080", "/health/liveness", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [timeout:%v]\n", dIn, dOut)

	//Output:
	//test: timeout(nil) -> [timeout:0s]
	//test: timeout(req) -> [duration:-1ns] [timeout:0s]
	//test: timeout(req) -> [duration:100ms] [timeout:100ms]

}

func ExampleTimeout_Deadline() {
	dIn := time.Millisecond * 200
	deadline := time.Millisecond * 100
	uri := "/search?q=golang"
	rsc := NewPrimaryResource("http://localhost:8080", "/health/liveness", dIn, httpCall)

	ctx, cancel := context.WithTimeout(context.Background(), deadline)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	dOut := rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	dIn = time.Millisecond * 100
	rsc = NewPrimaryResource("http://localhost:8080", "/health/liveness", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	dIn = time.Millisecond * 50
	rsc = NewPrimaryResource("http://localhost:8080", "/health/liveness", dIn, httpCall)
	dOut = rsc.timeout(req)
	fmt.Printf("test: timeout(req) -> [duration:%v] [deadline:%v] [timeout:%v]\n", dIn, deadline, dOut)

	//Output:
	//test: timeout(req) -> [duration:200ms] [deadline:100ms] [timeout:-100ms]
	//test: timeout(req) -> [duration:100ms] [deadline:100ms] [timeout:-100ms]
	//test: timeout(req) -> [duration:50ms] [deadline:100ms] [timeout:50ms]

}
