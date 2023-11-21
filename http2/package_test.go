package http2

import (
	"fmt"
	"reflect"
)

func Example_PkgUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	//pkgPath2 := runtime.PathFromUri(pkgUri2)

	fmt.Printf("test: PkgPath = \"%v\"\n", pkgUri2)
	//fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath2)

	//Output:
	//test: PkgPath = "github.com/advanced-go/core/http2"

}
