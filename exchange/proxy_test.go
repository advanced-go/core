package exchange

import (
	"fmt"
	"net/http"
)

func appHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func ExampleProxy_Add() {
	p := NewProxy()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	status := p.Register("", nil)
	fmt.Printf("test: Register(\"\") -> [status:%v]\n", status)

	status = p.Register(path, nil)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	status = p.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	status = p.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	status = p.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	//Output:
	//test: Register("") -> [status:Invalid Argument [invalid argument: path is empty]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler is nil: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler already exists: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK]

}

func ExampleProxy_Get() {
	p := NewProxy()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	_, status := p.Lookup("")
	fmt.Printf("test: Lookup(\"\") -> [status:%v]\n", status)

	_, status = p.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v]\n", path, status)

	status = p.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	handler, status1 := p.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", path, status1, handler != nil)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	status = p.Register(path, appHttpHandler)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)
	handler, status1 = p.Lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", path, status1, handler != nil)

	//Output:
	//test: Lookup("") -> [status:Invalid Argument [invalid argument: path is invalid: []]]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler does not exist: [github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK] [handler:true]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK]
	//test: Lookup(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK] [handler:true]

}
