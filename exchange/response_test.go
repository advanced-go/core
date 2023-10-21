package exchange

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/url"
)

func Example_ReadResponse() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	u, _ := url.Parse(s)
	//req, err := http.NewRequest("GET", s, nil)
	//fmt.Printf("test: http.NewRequest() -> [err:%v]\n", err)

	resp, err0 := ReadResponse(u)
	fmt.Printf("test: ReadResponse(%v) -> [err:%v] [status:%v]\n", s, err0, resp.StatusCode)

	buf, status := ReadAll[runtime.DebugError](resp.Body)
	fmt.Printf("test: ReadAll() -> [status:%v] %v", status.Code(), string(buf))

	s = string(buf)

	//Output:
	//test: ReadResponse(file://[cwd]/httptest/resource/http/html-response.html) -> [err:<nil>] [status:200]
	//test: ReadAll() -> [status:OK] <html>
	//<body>
	//<h1>Hello, World!</h1>
	//</body>
	//</html>

}
