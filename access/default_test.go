package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func Example_Formatter() {
	EnableTestLogger()
	start := time.Now().UTC()
	SetOrigin(Origin{Region: "us", Zone: "west", SubZone: "dc1", App: "search-app", InstanceId: "123456789"})

	req, err := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	req.Header.Add(runtime.XRelatesTo, "your-id")
	fmt.Printf("test: NewRequest() -> [err:%v] [req:%v]\n", err, req != nil)
	resp := http.Response{StatusCode: http.StatusOK}
	time.Sleep(time.Millisecond * 500)
	logTest(EgressTraffic, start, time.Since(start), req, &resp, "google-search", "secondary", -1, "")

	fmt.Printf("test: LogURI() -> %v\n", "success")

	//Output:
	//test: NewRequest() -> [err:<nil>] [req:true]
	//test: LogURI() -> success

}

func Example_Formatter_Urn() {
	start := time.Now().UTC()

	req, err := http.NewRequest("select", "github.com/advanced-go/example-domain/activity:entry", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	req.Header.Add(runtime.XRelatesTo, "fmtlog testing")
	fmt.Printf("test: NewRequest() -> [err:%v] [req:%v]\n", err, req != nil)
	resp := http.Response{StatusCode: http.StatusOK}
	logTest(InternalTraffic, start, time.Since(start), req, &resp, "route", "primary", -1, "")
	fmt.Printf("test: LogURN() -> %v\n", "success")

	//Output:
	//test: NewRequest() -> [err:<nil>] [req:true]
	//test: LogURN() -> success

}

func logTest(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, routeName, routeTo string, threshold int, thresholdFlags string) {
	Log(traffic, start, duration, req, resp, routeName, routeTo, threshold, thresholdFlags)
}
