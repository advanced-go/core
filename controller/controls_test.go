package controller

import (
	"fmt"
)

func ExampleControls_Add() {
	ctrl := NewController("test-route", nil, nil)
	p := NewControls()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	status := p.register("", nil)
	fmt.Printf("test: Register(\"\") -> [status:%v]\n", status)

	status = p.register(path, nil)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	status = p.register(path, ctrl)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	status = p.register(path, ctrl)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	status = p.register(path, ctrl)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	//Output:
	//test: Register("") -> [status:Invalid Argument [invalid argument: path is empty]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: Controller is nil: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: Controller already exists: [http://localhost:8080/github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK]

}

func ExampleControls_Get() {
	ctrl := NewController("test-route", nil, nil)
	p := NewControls()
	path := "http://localhost:8080/github.com/advanced-go/example-domain/activity"

	_, status := p.lookup("")
	fmt.Printf("test: Lookup(\"\") -> [status:%v]\n", status)

	_, status = p.lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v]\n", path, status)

	status = p.register(path, ctrl)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)

	handler, status1 := p.lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", path, status1, handler != nil)

	path = "http://localhost:8080/github/advanced-go/example-domain/activity"
	status = p.register(path, ctrl)
	fmt.Printf("test: Register(%v) -> [status:%v]\n", path, status)
	handler, status1 = p.lookup(path)
	fmt.Printf("test: Lookup(%v) -> [status:%v] [handler:%v]\n", path, status1, handler != nil)

	//Output:
	//test: Lookup("") -> [status:Invalid Argument [invalid argument: path is invalid: []]]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:Invalid Argument [invalid argument: Controller does not exist: [github.com/advanced-go/example-domain/activity]]]
	//test: Register(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK]
	//test: Lookup(http://localhost:8080/github.com/advanced-go/example-domain/activity) -> [status:OK] [handler:true]
	//test: Register(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK]
	//test: Lookup(http://localhost:8080/github/advanced-go/example-domain/activity) -> [status:OK] [handler:true]

}
