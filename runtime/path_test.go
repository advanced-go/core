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

func Example_PathFromUri() {
	s := "github.com/advanced-go/core/runtime"
	p := PathFromUri2(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	s = "github.comadvanced-gocoreruntime"
	p = PathFromUri2(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	//Output:
	//test: PathFromUri(github.com/advanced-go/core/runtime) /advanced-go/core/runtime
	//test: PathFromUri(github.comadvanced-gocoreruntime) [uri invalid]

}
