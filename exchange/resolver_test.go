package exchange

import "fmt"

func echo(key string) string {
	return key
}

func empty(key string) string {
	return ""
}

func Example_Resolver() {
	r := NewResolver("http://localhost:8080", echo)

	key := "test"
	val := r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	key = ""
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	key = "/"
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	r = NewResolver("http://localhost:8080", empty)
	key = "/"
	val = r.Resolve(key)
	fmt.Printf("test: Resolver(\"%v\") -> %v\n", key, val)

	//Output:
	//test: Resolver("test") -> test
	//test: Resolver("") ->
	//test: Resolver("/") -> /
	//test: Resolver("/") -> http://localhost:8080/

}

/*
func Example_resolve() {
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

func Example_addResolver() {
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


*/
