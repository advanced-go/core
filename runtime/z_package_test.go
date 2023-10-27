package runtime

import (
	"fmt"
)

func Example_PackageUri() {
	fmt.Printf("test: PkgUri -> %v\n", PkgUri)

	//Output:
	//test: PkgUri -> github.com/go-ai-agent/core/runtime

}

func Example_PathFromUri() {
	s := "github.com/go-ai-agent/core/runtime"
	p := PathFromUri(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	s = "github.comgo-ai-agentcoreruntime"
	p = PathFromUri(s)
	fmt.Printf("test: PathFromUri(%v) %v\n", s, p)

	//Output:
	//test: PathFromUri(github.com/go-ai-agent/core/runtime) /go-ai-agent/core/runtime
	//test: PathFromUri(github.comgo-ai-agentcoreruntime) [uri invalid]

}
