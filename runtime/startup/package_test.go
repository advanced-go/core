package startup

import (
	"fmt"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUri -> %v\n", PkgUri)

	//Output:
	//test: PkgUri -> github.com/advanced-go/core/runtime/startup

}

func Example_newStatusRequest() {
	req := newStatusRequest()
	fmt.Printf("test: newStatusRequest() -> %v\n", req.URL.Path)

	//Output:
	//test: newStatusRequest() -> /startup/status

}
