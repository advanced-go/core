package exchange

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
)

const (
	testLocation = "github.com/advanced-go/core/exchange:ExampleGet"
)

func ExampleGet() {
	runtime.SetOutputFormatter()
	var e runtime.Output
	r, status := Get("", nil)

	fmt.Printf("test: Get(\"\") -> [resp:%v] [status:%v]\n", r.Status, status)
	e.Handle(status, "123-456", testLocation)

	//Output:
	//test: Get("") -> [resp:Internal Error] [status:Internal Error [Get "": unsupported protocol scheme ""]]

}
