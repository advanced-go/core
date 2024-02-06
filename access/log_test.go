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
	s := DefaultFormatter(&Origin{Region: "us", Zone: "zone", App: "ai-agent"}, EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: 200, Status: "OK"}, "route", "primary", -1, "")

	fmt.Printf("test: fmtlog() -> %v\n", s)

	//Output:

}

func Example_NewRequest() {
	h := make(http.Header)
	h.Add("key-1", "value-1")
	h.Add("key-2", "value-2")
	h.Add(XRequestId, "123-456")

	r := NewRequest(h, http.MethodGet, "https://www.google.com/search?q=golang")
	fmt.Printf("test: NewRequest() -> [method:%v] [uri:%v] [header:%v]\n", r.Method, r.URL.String(), r.Header)

	r = NewRequest(h, http.MethodPatch, failsafeUri)
	fmt.Printf("test: NewRequest() -> [method:%v] [uri:%v] [header:%v]\n", r.Method, r.URL.String(), r.Header)

	//Output:
	//test: NewRequest() -> [method:GET] [uri:https://www.google.com/search?q=golang] [header:map[Key-1:[value-1] Key-2:[value-2] X-Request-Id:[123-456]]]
	//test: NewRequest() -> [method:PATCH] [uri:https://invalid-uri.com] [header:map[Key-1:[value-1] Key-2:[value-2] X-Request-Id:[123-456]]]

}

type testStatus struct {
	code int
}

func (ts *testStatus) StatusCode() int {
	return ts.code
}

func (ts *testStatus) String() string {
	return fmt.Sprintf("%v", ts.code)
}

func ExampleLogDeferred() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
	//test: LogDeferred() -> Timeout

}

func loggingTest() *runtime.Status {
	var status *runtime.Status

	h := make(http.Header)
	h.Add(XRequestId, XRequestId)
	h.Add(XRelatesTo, XRelatesTo)
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search",
		"us.west", -1, "flags", StatusCode(&status))()
	status = runtime.NewStatus(http.StatusGatewayTimeout)
	time.Sleep(time.Millisecond * 500)
	return status
}

/*
func ExampleLogDeferred_Test1() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest1()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
	//test: LogDeferred() -> 504
}

// default status variable
func loggingTest1() (statusCode int) {
	h := make(http.Header)
	h.Add(XRequestId, XRequestId)
	h.Add(XRelatesTo, XRelatesTo)
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search", "us.west", -1, "flags", StatusCode(&statusCode))()
	statusCode = http.StatusGatewayTimeout
	time.Sleep(time.Millisecond * 500)
	return
}

func ExampleLogDeferred_Test2() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest2()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
	//test: LogDeferred() -> 503

}

// non-default status variable
func loggingTest2() int {
	var statusCode int
	h := make(http.Header)
	h.Add(XRequestId, XRequestId)
	h.Add(XRelatesTo, XRelatesTo)
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search", "us.west", -1, "flags", StatusCode(&statusCode))()
	statusCode = http.StatusServiceUnavailable
	time.Sleep(time.Millisecond * 500)
	return statusCode
}

type StatusCode2 interface {
	StatusCode() int
}

func NewStatusCode(t any) func() int {
	return func() int {
		if t == nil {
			return http.StatusOK
		}
		fmt.Printf("reflect.TypeOf() -> %v\n", reflect.TypeOf(t).String())
		//reflect.Interface
		if i, ok := t.(*StatusCode2); ok {
			return (*(i)).StatusCode()
		}
		return 0
		//return (*(s)).StatusCode()
	}
}

*/
