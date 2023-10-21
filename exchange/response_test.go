package exchange

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

func _Example_createFileResponse() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	req, err := http.NewRequest("GET", s, nil)

	fmt.Printf("test: http.NewRequest() -> [err:%v]\n", err)

	resp, err0 := createFileResponse(req)
	fmt.Printf("test: createFileResponse() -> [err:%v] [status:%v]\n", err0, resp.StatusCode)

	buf, status := ReadAll[runtime.DebugError](resp.Body)
	fmt.Printf("test: ReadAll() -> [err:%v] %v", status.Code(), string(buf))

	s = string(buf)

	//Output:
	//test: http.NewRequest() -> [err:<nil>]
	//test: createFileResponse() -> [err:<nil>] [status:200]
	//test: ReadAll() -> [err:<nil>] <html>
	//<body>
	//<h1>Hello, World!</h1>
	//</body>
	//</html>

}
