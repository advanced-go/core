package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"reflect"
)

func Example_PkgUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath2 := runtime.PathFromUri(pkgUri2)

	fmt.Printf("test: PkgUri  = \"%v\"\n", pkgUri2)
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath2)

	//Output:
	//test: PkgUri  = "github.com/advanced-go/core/http2"
	//test: PkgPath = "/advanced-go/core/http2"

}
