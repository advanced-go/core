package runtime

import (
	"fmt"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUrl -> %v\n", PkgUrl)
	fmt.Printf("test: PkgUri -> %v\n", PkgUri)
	fmt.Printf("test: PkgPath -> %v\n", PkgUrl.Path)

	//Output:
	//test: PkgUrl -> file://github.com/go-ai-agent/core/runtime
	//test: PkgUri -> github.com/go-ai-agent/core/runtime
	//test: PkgPath -> /go-ai-agent/core/runtime
	
}
