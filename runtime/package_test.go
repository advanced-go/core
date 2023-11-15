package runtime

import (
	"fmt"
	"reflect"
)

func Example_PackageUri() {
	pkgUri2 := reflect.TypeOf(any(pkg{})).PkgPath()
	pkgPath2 := PathFromUri(pkgUri2)

	fmt.Printf("test: PkgUri  = \"%v\"\n", pkgUri2)
	fmt.Printf("test: PkgPath = \"%v\"\n", pkgPath2)

	//Output:
	//test: PkgUri  = "github.com/advanced-go/core/runtime"
	//test: PkgPath = "/advanced-go/core/runtime"

}

func Example_PathFromUri() {
	s := "github.com/advanced-go/core/runtime"
	p := PathFromUri(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	s = "github.comadvanced-gocoreruntime"
	p = PathFromUri(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	//Output:
	//test: PathFromUri(github.com/advanced-go/core/runtime) /advanced-go/core/runtime
	//test: PathFromUri(github.comadvanced-gocoreruntime) [uri invalid]

}
