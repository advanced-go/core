package runtime

import (
	"fmt"
	"net/http"
)

func appHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func Example_HandlerMap_Add() {
	m := NewHandlerMap()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	status := m.AddHandler("", nil)
	fmt.Printf("test: AddHandler(\"\") -> [status:%v]\n", status)

	status = m.AddHandler(path, nil)
	fmt.Printf("test: AddHandler(%v) -> [status:%v]\n", path, status)

	status = m.AddHandler(path, appHttpHandler)
	fmt.Printf("test: AddHandler(%v) -> [status:%v]\n", path, status)

	status = m.AddHandler(path, appHttpHandler)
	fmt.Printf("test: AddHandler(%v) -> [status:%v]\n", path, status)

	//Output:
	//test: AddHandler("") -> [status:Invalid Argument [invalid argument: path is empty]]
	//test: AddHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler is nil: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]
	//test: AddHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: AddHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler already exists: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]

}

func Example_HandlerMap_Get() {
	m := NewHandlerMap()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	_, status := m.GetHandler("")
	fmt.Printf("test: GetHandler(\"\") -> [status:%v]\n", status)

	_, status = m.GetHandler(path)
	fmt.Printf("test: GetHandler(%v) -> [status:%v]\n", path, status)

	status = m.AddHandler(path, appHttpHandler)
	fmt.Printf("test: AddHandler(%v) -> [status:%v]\n", path, status)

	handler, status1 := m.GetHandler(path)
	fmt.Printf("test: GetHandler(%v) -> [status:%v] [handler:%v]\n", path, status1, handler != nil)

	//Output:
	//test: GetHandler("") -> [status:Invalid Argument [invalid argument: path is invalid: []]]
	//test: GetHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: HTTP handler does not exist: [github.com/advanced-go/example-domain/activity]]]
	//test: AddHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: GetHandler(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK] [handler:true]

}
