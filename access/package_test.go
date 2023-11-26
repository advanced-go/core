package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"reflect"
)

func Example_PackageUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgUri  = \"%v\"\n", pkgUri2)

	//Output:
	//test: PkgUri  = "github.com/advanced-go/core/access"
	
}

func Example_LogAccess() {
	// w := WrapDo[defaultLogFn](nil,nil,nil)
}

func ExampleNewStatusCodeFn() {
	var status runtime.Status

	fn := NewStatusCodeClosure(&status)
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
	fmt.Printf("test: NewStatusCode(&status) -> [statusCode:%v]\n", fn())

	//Output:
	//test: NewStatusCode(&status) -> [statusCode:4]

}
