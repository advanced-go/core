package http2test

import (
	"fmt"
	"net/http"
)

func Example_DoT() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := DoT[[]byte](req)
	fmt.Printf("test: DoT[[]byte](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: DoT[[]byte](req) -> [status:OK] [buf:false] [resp:true]

}
