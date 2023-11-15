package http2test

import (
	"fmt"
	"github.com/advanced-go/core/io2"
	"net/url"
)

func Example_ReadResponse() {
	s := "file://[cwd]/http2test/resource/html-response.txt"
	u, _ := url.Parse(s)
	//req, err := http.NewRequest("GET", s, nil)
	//fmt.Printf("test: http.NewRequest() -> [err:%v]\n", err)

	resp, err0 := ReadResponse(u)
	fmt.Printf("test: ReadResponse(%v) -> [err:%v] [status:%v]\n", s, err0, resp.StatusCode)

	buf, status := io2.ReadAll(resp.Body)
	fmt.Printf("test: ReadAll() -> [status:%v] %v", status, string(buf))

	s = string(buf)

	//Output:
	//test: ReadResponse(file://[cwd]/http2test/resource/html-response.txt) -> [err:<nil>] [status:200]
	//test: ReadAll() -> [status:OK] <html>
	//<body>
	//<h1>Hello, World!</h1>
	//</body>
	//</html>

}
