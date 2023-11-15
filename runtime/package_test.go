package runtime

import (
	"fmt"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUri -> %v\n", PkgUri)

	//Output:
	//test: PkgUri -> github.com/advanced-go/core/runtime

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
