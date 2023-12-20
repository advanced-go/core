package io

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/url"
)

func readAll(body io.ReadCloser) ([]byte, runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, ":ReadAll", err)
	}
	return buf, runtime.StatusOK()
}

func Example_ReadResponse() {
	s := "file://[cwd]/iotest/html-response.txt"
	u, _ := url.Parse(s)
	//req, err := http.NewRequest("GET", s, nil)
	//fmt.Printf("test: http.NewRequest() -> [err:%v]\n", err)

	resp, err0 := ReadResponse(u)
	fmt.Printf("test: ReadResponse(%v) -> [err:%v] [status:%v]\n", s, err0, resp.StatusCode)

	buf, status := readAll(resp.Body)
	fmt.Printf("test: readAll() -> [status:%v] %v", status, string(buf))

	s = string(buf)

	//Output:
	//test: ReadResponse(file://[cwd]/iotest/html-response.txt) -> [err:<nil>] [status:200]
	//test: ReadAll() -> [status:OK] <html>
	//<body>
	//<h1>Hello, World!</h1>
	//</body>
	//</html>

}
