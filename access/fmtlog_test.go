package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

func Example_FmtLog() {
	start := time.Now().UTC()

	req, err := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	fmt.Printf("test: NewRequest() -> [err:%v] %v\n", err, req)
	resp := http.Response{StatusCode: http.StatusOK}
	s := fmtLog(EgressTraffic, start, time.Since(start), req, &resp, "route", -1, "")
	fmt.Printf("test: fmtLog() -> %v\n", s)

	//Output:

}

func Example_FmtLog_Urn() {
	start := time.Now().UTC()

	req, err := http.NewRequest("select", "github.com/advanced-go/example-domain/activity:entry", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	req.Header.Add("RelatesTo", "fmtlog testing")
	fmt.Printf("test: NewRequest() -> [err:%v] %v\n", err, req)
	resp := http.Response{StatusCode: http.StatusOK}
	s := fmtLog(InternalTraffic, start, time.Since(start), req, &resp, "route", -1, "")
	fmt.Printf("test: fmtLog() -> %v\n", s)

	//Output:

}
