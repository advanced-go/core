package runtime

import (
	"fmt"
	"reflect"
)

func Example_PackageUri() {
	pkgPath := reflect.TypeOf(any(pkg{})).PkgPath()
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath)
	//fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath2)

	//Output:
	//test: PkgPath = "github.com/advanced-go/core/runtime"

}
