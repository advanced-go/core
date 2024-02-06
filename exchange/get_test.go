package exchange

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

const (
	testLocation = "github/advanced-go/core/exchange:ExampleGet"
)

func ExampleGet() {
	runtime.SetOutputFormatter()
	var e runtime.Output
	r, status := Get(nil, "", nil)

	e.Handle(status, "123-456", testLocation)
	fmt.Printf("test: Get(\"\") -> [resp:%v] [status:%v]\n", r.Status, status)

	//Output:
	//test: Get("") -> [resp:Internal Error] [status:Bad Request [error: URI is empty]]

}
