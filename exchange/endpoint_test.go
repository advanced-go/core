package exchange

import (
	"fmt"
)

func Example_RegisterEndpoint() {
	path := ""

	status := RegisterEndpoint(path, nil)
	fmt.Printf("test: RegisterEndpoint(\"\") -> [status:%v]\n", status)

	path = ""

	status = RegisterEndpoint(path, nil)
	fmt.Printf("test: RegisterEndpoint(\"\") -> [status:%v]\n", status)

	//Output:

}
