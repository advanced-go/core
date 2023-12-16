package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func Example_Formatter() {
	start := time.Now().UTC()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})

	req, err := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	req.Header.Add(runtime.XRelatesTo, "your-id")
	fmt.Printf("test: NewRequest() -> [err:%v] %v\n", err, req)
	resp := http.Response{StatusCode: http.StatusOK}
	s := defaultFormatter(origin, EgressTraffic, start, time.Since(start), req, &resp, "google-search", "secondary", -1, "")
	fmt.Printf("test: formatter() -> %v\n", s)

	//Output:

}

func Example_Formatter_Urn() {
	start := time.Now().UTC()

	req, err := http.NewRequest("select", "github.com/advanced-go/example-domain/activity:entry", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	req.Header.Add(runtime.XRelatesTo, "fmtlog testing")
	fmt.Printf("test: NewRequest() -> [err:%v] %v\n", err, req)
	resp := http.Response{StatusCode: http.StatusOK}
	s := defaultFormatter(origin, InternalTraffic, start, time.Since(start), req, &resp, "route", "primary", -1, "")
	fmt.Printf("test: fmtLog() -> %v\n", s)

	//Output:

}
