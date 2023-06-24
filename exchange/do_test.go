package exchange

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

func ExampleDoT() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := DoT[runtime.DebugError, Default, []byte](req)
	fmt.Printf("test: Do[DebugError,[]byte,DefaultExchange](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: Do[DebugError,[]byte,DefaultExchange](req) -> [status:OK] [buf:true] [resp:true]
}
