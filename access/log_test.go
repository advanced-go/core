package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func _Example_LogAccess() {
	start := time.Now().UTC()
	r, _ := http.NewRequest("PUT", "/github.com/advanced-go/example-domain/activity:entry", nil)
	r.Host = "localhost:8080"
	s := defaultFormatter(Origin{Region: "us", Zone: "zone", App: "ai-agent"}, EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: 200, Status: "OK"}, "route", "primary", -1, "")

	fmt.Printf("test: fmtlog() -> %v\n", s)

	//Output:

}

func Example_NewRequest() {
	h := make(http.Header)
	h.Add("key-1", "value-1")
	h.Add("key-2", "value-2")
	h.Add(runtime.XRequestId, "123-456")

	r := NewRequest(h, http.MethodGet, "https://www.google.com/search?q=golang")
	fmt.Printf("test: NewRequest() -> [method:%v] [uri:%v] [header:%v]\n", r.Method, r.URL.String(), r.Header)

	r = NewRequest(h, http.MethodPatch, failsafeUri)
	fmt.Printf("test: NewRequest() -> [method:%v] [uri:%v] [header:%v]\n", r.Method, r.URL.String(), r.Header)

	//Output:
	//test: NewRequest() -> [method:GET] [uri:https://www.google.com/search?q=golang] [header:map[Key-1:[value-1] Key-2:[value-2] X-Request-Id:[123-456]]]
	//test: NewRequest() -> [method:PATCH] [uri:https://invalid-uri.com] [header:map[Key-1:[value-1] Key-2:[value-2] X-Request-Id:[123-456]]]

}

func ExampleLogDeferred_Test1() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest1()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
	//test: LogDeferred() -> Timeout
}

// default status variable
func loggingTest1() (status runtime.Status) {
	h := make(http.Header)
	h.Add(runtime.XRequestId, runtime.XRequestId)
	h.Add(runtime.XRelatesTo, runtime.XRelatesTo)
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search", "us.west", -1, "flags", &status)()
	status = runtime.NewStatus(http.StatusGatewayTimeout)
	time.Sleep(time.Millisecond * 500)
	return
}

func ExampleLogDeferred_Test2() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest2()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
	//test: LogDeferred() -> Service Unavailable

}

// non-default status variable
func loggingTest2() runtime.Status {
	var status runtime.Status
	h := make(http.Header)
	h.Add(runtime.XRequestId, runtime.XRequestId)
	h.Add(runtime.XRelatesTo, runtime.XRelatesTo)
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search", "us.west", -1, "flags", &status)()
	status = runtime.NewStatus(http.StatusServiceUnavailable)
	time.Sleep(time.Millisecond * 500)
	return status
}
