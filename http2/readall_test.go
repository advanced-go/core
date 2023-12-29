package http2

import (
	"fmt"
	"reflect"
)

func Example_ReadAll() {
	var cnt = 5

	_, status := ReadAll(cnt)
	fmt.Printf("test: ReadAll(%v) [status:%v]\n", reflect.TypeOf(cnt), status)

	//Output:
	//test: ReadAll(int) [status:Invalid Argument [int]]

}
