package access

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"reflect"
)

func Example_PackageUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath2 := runtime.PathFromUri(pkgUri2)

	fmt.Printf("test: PkgUri  = \"%v\"\n", pkgUri2)
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath2)

	//Output:
	//test: PkgUri  = "github.com/advanced-go/core/access"
	//test: PkgPath = "/advanced-go/core/access"

}

func Example_LogAccess() {
	// w := WrapDo[defaultLogFn](nil,nil,nil)
}