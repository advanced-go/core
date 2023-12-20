package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"reflect"
	"time"
)

func Example_PackageUri() {
	pkgPath := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgPath  = \"%v\"\n", pkgPath)

	//Output:
	//test: PkgPath  = "github.com/advanced-go/core/access"

}

func Example_LogAccess() {
	start := time.Now().UTC()
	r, _ := http.NewRequest("PUT", "/github.com/advanced-go/example-domain/activity:entry", nil)
	r.Host = "localhost:8080"
	s := defaultFormatter(Origin{Region: "us", Zone: "zone", App: "ai-agent"}, EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: 200}, "route", "primary", -1, "")

	fmt.Printf("test: fmtlog() -> %v\n", s)

	//Output:

}

func Example_NewRequest() {
	h := make(http.Header)
	h.Add("key-1", "value-1")
	h.Add("key-2", "value-2")
	h.Add(runtime.XRequestId, "123-456")

	r := NewRequest(h, "GET", "https://www.google.com/search?q=golang")
	fmt.Printf("test: NewRequest() -> [method:%v] [uri:%v] [header:%v]\n", r.Method, r.URL.String(), r.Header)

	//Output:
	//test: NewRequest() -> [method:GET] [uri:https://www.google.com/search?q=golang] [header:map[Key-1:[value-1] Key-2:[value-2] X-Request-Id:[123-456]]]

}

func ExampleLogDeferred_Test1() {
	EnableTestLogger()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})
	status := loggingTest1()
	fmt.Printf("test: LogDeferred() -> %v\n", status)

	//Output:
}

// default status variable
func loggingTest1() (status runtime.Status) {
	h := make(http.Header)
	h.Add(runtime.XRequestId, "x-request-id")
	h.Add(runtime.XRelatesTo, "x-relates-to")
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
}

// non-default status variable
func loggingTest2() runtime.Status {
	var status runtime.Status
	h := make(http.Header)
	h.Add(runtime.XRequestId, "x-request-id")
	h.Add(runtime.XRelatesTo, "x-relates-to")
	defer LogDeferred(EgressTraffic, NewRequest(h, http.MethodGet, "https://www.google.com/search?q=test"), "search", "us.west", -1, "flags", &status)()
	status = runtime.NewStatus(http.StatusServiceUnavailable)
	time.Sleep(time.Millisecond * 500)
	return status
}
