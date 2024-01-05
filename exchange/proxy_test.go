package exchange

import (
	"fmt"
)

func Example_RegisterHandler() {
	path := ""
	status := RegisterHandler(path, nil)
	fmt.Printf("test: RegisterEndpoint(\"\") -> [status:%v]\n", status)

	//Output:
	//test: RegisterEndpoint("") -> [status:Invalid Argument [invalid argument: path is empty]]

}
