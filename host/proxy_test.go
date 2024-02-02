package host

import (
	"fmt"
)

/*
func proxyTestHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

*/

func ExampleProxy_Add() {
	proxy := NewProxy()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	err := proxy.Register("", nil)
	fmt.Printf("test: Register(\"\") -> [err:%v]\n", err)

	err = proxy.Register(path, nil)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)

	err = proxy.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)

	err = proxy.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	err = proxy.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)

	//Output:
	//test: Register("") -> [err:error: proxy.Register() path is empty]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [err:error: proxy.Register() HTTP handler is nil: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [err:<nil>]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [err:error: proxy.Register() HTTP handler already exists: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [err:<nil>]

}

func ExampleProxy_Get() {
	proxy := NewProxy()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	p := proxy.Lookup("")
	fmt.Printf("test: Lookup(\"\") -> [proxy:%v]\n", p != nil)

	p = proxy.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [proxy:%v]\n", path, p != nil)

	err := proxy.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)

	handler := proxy.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [handler:%v]\n", path, handler != nil)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	err = proxy.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [err:%v]\n", path, err)
	handler = proxy.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [handler:%v]\n", path, handler != nil)

	//Output:
	//test: Lookup("") -> [proxy:false]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [proxy:false]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [err:<nil>]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [handler:true]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [err:<nil>]
	//test: Lookup(http://localhost:8080/github/advanced-go/example-domain/activity) -> [handler:true]

}
