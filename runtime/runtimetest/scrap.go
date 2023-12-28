package runtimetest

import (
	"fmt"
	"strings"
)

// PathFromUri2 - return a path from a scheme less uri
func PathFromUri2(rawUri string) string {
	i := strings.Index(rawUri, "/")
	if i < 0 {
		return "[uri invalid]"
	}
	return rawUri[i:]
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
