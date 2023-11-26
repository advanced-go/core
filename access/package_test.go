package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"reflect"
	"time"
)

func Example_PackageUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgUri  = \"%v\"\n", pkgUri2)

	//Output:
	//test: PkgUri  = "github.com/advanced-go/core/access"

}

func ExampleNewStatusCodeFn() {
	var status runtime.Status

	fn := NewStatusCodeClosure(&status)
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
	fmt.Printf("test: NewStatusCode(&status) -> [statusCode:%v]\n", fn())

	//Output:
	//test: NewStatusCode(&status) -> [statusCode:4]

}

func Example_LogAccess() {
	start := time.Now().UTC()
	r, _ := http.NewRequest("PUT", "/github.com/advanced-go/example-domain/activity:entry", nil)
	r.Host = "localhost:8080"
	s := fmtLog(EgressTraffic, start, time.Since(start), r, &http.Response{StatusCode: 200}, -1, "")

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
