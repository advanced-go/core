package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [Invalid Argument [invalid argument : request is nil]]

}

func ExampleDo_GatewayTimeout() {
	status1 := runtime.NewStatus(http.StatusGatewayTimeout)
	req, _ := http.NewRequestWithContext(NewStatusContext(nil, status1), http.MethodGet, "file://[cwd]/http2test/resource/http-503.txt", nil)
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [status:%v] [body:%v]\n",
		resp != nil, status.Code(), status, resp.Body != nil)

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:504] [status:Timeout] [body:false]

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

func ExampleDo_ServiceUnavailable_ContentLocation() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add(ContentLocation, "file://[cwd]/http2test/resource/http-503.txt")
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [errs:%v] [content-type:%v] [body:%v]\n",
		resp != nil, status.Code(), status.Errors(), resp.Header.Get("content-type"), resp.Body != nil)

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [errs:[]] [content-type:text/html] [body:true]

}

func Example_DoT() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := DoT[[]byte](req)
	fmt.Printf("test: DoT[[]byte](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: DoT[[]byte](req) -> [status:OK] [buf:true] [resp:true]

}
