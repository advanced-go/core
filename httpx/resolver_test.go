package httpx

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

func Example_Resolve() {
	var s = ""
	url := Resolve(s)

	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "http://"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "/test/resource?env=dev&cust=1"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	s = "https://www.google.com/search?q=testing"
	url = Resolve(s)
	fmt.Printf("test: Resolve(%v) -> [%v]\n", s, url)

	//Output:
	//test: Resolve() -> []
	//test: Resolve(http://) -> [http://]
	//test: Resolve(/test/resource?env=dev&cust=1) -> [http://localhost:8080/test/resource?env=dev&cust=1]
	//test: Resolve(https://www.google.com/search?q=testing) -> [https://www.google.com/search?q=testing]

}

func Example_AddResolver() {
	pattern := "/endpoint/resource"

	uri := Resolve(pattern)
	fmt.Printf("test: Resolve(%v) -> %v\n", pattern, uri)

	AddResolver(func(s string) string {
		if s == pattern {
			return "https://github.com/acccount/go-ai-agent/core"
		}
		return ""
	})

	uri = Resolve("invalid")
	fmt.Printf("test: Resolve(%v) -> %v\n", pattern, uri)

	uri = Resolve(pattern)
	fmt.Printf("test: Resolve(%v) -> %v\n", pattern, uri)

	pattern2 := "/endpoint/resource2"
	AddResolver(func(s string) string {
		if s == pattern2 {
			return "https://gitlab.com/entry/idiomatic-go"
		}
		return ""
	})

	uri = Resolve(pattern2)
	fmt.Printf("test: Resolve(%v) -> %v\n", pattern2, uri)

	//Output:
	//test: Resolve(/endpoint/resource) -> http://localhost:8080/endpoint/resource
	//test: Resolve(/endpoint/resource) -> invalid
	//test: Resolve(/endpoint/resource) -> https://github.com/acccount/go-ai-agent/core
	//test: Resolve(/endpoint/resource2) -> https://gitlab.com/entry/idiomatic-go

}

func Example_AddResolver_Fail() {
	runtime.SetProdEnvironment()
	pattern := "/endpoint/resource"

	AddResolver(func(s string) string {
		if s == pattern {
			return "https://github.com/acccount/go-ai-agent/core"
		}
		return ""
	})

	fmt.Printf("test: AddResolver(%v) -> [err:%v]\n", pattern, nil)

	//Output:
	//test: AddResolver(/endpoint/resource) -> [err:<nil>]

}
