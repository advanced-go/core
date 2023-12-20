package http2test

import (
	"fmt"
	"net/http"
)

func ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [Invalid Argument [invalid argument : request is nil]]

}

func ExampleDo_ServiceUnavailable_Uri() {
	req, _ := http.NewRequest(http.MethodGet, "file://[cwd]/http2test/resource/http-503.txt", nil)
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [errs:%v] [content-type:%v] [body:%v]\n",
		resp != nil, status.Code(), status.Errors(), resp.Header.Get("content-type"), resp.Body != nil)

	//defer resp.Body.Close()
	//buf, ioError := io.ReadAll(resp.Body)
	//fmt.Printf("test: ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [errs:[]] [content-type:text/html] [body:true]

}