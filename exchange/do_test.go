package exchange

import (
	"fmt"
	"github.com/go-sre/core/runtime"
	"net/http"
)

func ExampleDo() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := Do[runtime.DebugError2, Default, []byte](req)
	fmt.Printf("test: Do[DebugError,[]byte,DefaultExchange](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: Do[DebugError,[]byte,DefaultExchange](req) -> [status:OK] [buf:true] [resp:true]
}
