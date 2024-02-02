package messaging

import (
	"fmt"
	"net/http"
)

func ExampleNewStatus() {
	s := NewStatus(http.StatusTeapot)

	fmt.Printf("test: NewStatus() -> [status:%v] [ok:%v]\n", s, s.OK())

	//Output:
	//test: NewStatus() -> [status:418] [ok:false]
	
}
