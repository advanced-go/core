package exchange

import "fmt"

func ExampleGet() {
	r, status := Get("", nil)

	fmt.Printf("test: Get(\"\") -> [resp:%v] [status:%v]\n", r.Status, status)

	//Output:
	//test: Get("") -> [resp:Internal Error] [status:Internal Error [Get "": unsupported protocol scheme ""]]

}
