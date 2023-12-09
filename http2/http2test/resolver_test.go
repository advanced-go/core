package http2test

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

func Example_Resolve() {
	var s = ""
	url := resolve(s)

	fmt.Printf("test: resolve(%v) -> [%v]\n", s, url)

	s = "http://"
	url = resolve(s)
	fmt.Printf("test: resolve(%v) -> [%v]\n", s, url)

	s = "/test/resource?env=dev&cust=1"
	url = resolve(s)
	fmt.Printf("test: resolve(%v) -> [%v]\n", s, url)

	s = "https://www.google.com/search?q=testing"
	url = resolve(s)
	fmt.Printf("test: resolve(%v) -> [%v]\n", s, url)

	//Output:
	//test: resolve() -> []
	//test: resolve(http://) -> [http://]
	//test: resolve(/test/resource?env=dev&cust=1) -> [http://localhost:8080/test/resource?env=dev&cust=1]
	//test: resolve(https://www.google.com/search?q=testing) -> [https://www.google.com/search?q=testing]

}

func Example_AddResolver() {
	pattern := "/endpoint/resource"

	uri := resolve(pattern)
	fmt.Printf("test: resolve(%v) -> %v\n", pattern, uri)

	addResolver(func(s string) string {
		if s == pattern {
			return "https://github.com/acccount/go-ai-agent/core"
		}
		return ""
	})

	uri = resolve("invalid")
	fmt.Printf("test: resolve(%v) -> %v\n", pattern, uri)

	uri = resolve(pattern)
	fmt.Printf("test: resolve(%v) -> %v\n", pattern, uri)

	pattern2 := "/endpoint/resource2"
	addResolver(func(s string) string {
		if s == pattern2 {
			return "https://gitlab.com/entry/idiomatic-go"
		}
		return ""
	})

	uri = resolve(pattern2)
	fmt.Printf("test: resolve(%v) -> %v\n", pattern2, uri)

	//Output:
	//test: resolve(/endpoint/resource) -> http://localhost:8080/endpoint/resource
	//test: resolve(/endpoint/resource) -> invalid
	//test: resolve(/endpoint/resource) -> https://github.com/acccount/go-ai-agent/core
	//test: resolve(/endpoint/resource2) -> https://gitlab.com/entry/idiomatic-go

}

func Example_AddResolver_Fail() {
	runtime.SetProdEnvironment()
	pattern := "/endpoint/resource"

	addResolver(func(s string) string {
		if s == pattern {
			return "https://github.com/acccount/go-ai-agent/core"
		}
		return ""
	})

	fmt.Printf("test: addResolver(%v) -> [err:%v]\n", pattern, nil)

	//Output:
	//test: addResolver(/endpoint/resource) -> [err:<nil>]

}
