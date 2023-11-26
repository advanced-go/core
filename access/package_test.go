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
