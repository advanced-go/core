package exchange

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	testContent = "this is response write content"
	requestId   = "123-request-id"
	relatesTo   = "test-relates-to"
	statusCode  = http.StatusAccepted
)

func ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [Invalid Argument [invalid argument : request is nil]]

}

func ExampleDo_ServiceUnavailable_Uri() {
	req, _ := http.NewRequest(http.MethodGet, "file://[cwd]/exchangetest/http-503.txt", nil)
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [errs:%v] [content-type:%v] [body:%v]\n",
		resp != nil, status.Code, status.Error, resp.Header.Get("content-type"), resp.Body != nil)

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [errs:<nil>] [content-type:text/html] [body:true]

}

/*
func ExampleDo_ConnectivityError() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, status := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [errs:%v] [content-type:%v] [body:%v]\n",
		resp != nil, status.Code(), status.Errors(), resp.Header.Get("content-type"), resp.Body != nil)

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [errs:[]] [content-type:text/html] [body:true]

}


*/

func ExampleDo_Service_Unavailable() {
	s := "file://[cwd]/exchangetest/http-503.txt"
	req, _ := http.NewRequest("", s, nil)
	resp, status := Do(req)
	fmt.Printf("test: Do() -> [status-code:%v] [status:%v]\n", resp.StatusCode, status)

	//Output:
	//test: Do() -> [status-code:503] [status:Service Unavailable]

}

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(runtime.XRequestId, requestId)
	w.Header().Add(runtime.XRelatesTo, relatesTo)
	w.WriteHeader(statusCode)
	w.Write([]byte(testContent))
}

func ExampleDo_Proxy() {
	uri := "http://localhost:8080/github.com/advanced-go/core/exchange:Do"
	req, _ := http.NewRequest("", uri, nil)

	resp, status := Do(req)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp != nil, status)

	status = RegisterHandler(uri, testHandler)
	fmt.Printf("test: RegisterEndpoint() -> [status:%v]\n", status)

	resp, status = Do(req)
	fmt.Printf("test: Do() -> [resp:%v] [status:%v]\n", resp != nil, status)

	fmt.Printf("test: Do() -> [write-requestId:%v] [response-requestId:%v]\n", requestId, resp.Header.Get(runtime.XRequestId))
	fmt.Printf("test: Do() -> [write-relatesTo:%v] [response-relatesTo:%v]\n", relatesTo, resp.Header.Get(runtime.XRelatesTo))
	fmt.Printf("test: Do() -> [write-statusCode:%v] [response-statusCode:%v]\n", statusCode, resp.StatusCode)

	buf, _ := runtime.ReadAll(resp.Body, nil)
	fmt.Printf("test: Do() -> [write-content:%v] [response-content:%v]\n", testContent, string(buf))

	//Output:
	//test: Do() -> [resp:true] [status:Internal Error [Get "http://localhost:8080/github.com/advanced-go/core/exchange:Do": dial tcp [::1]:8080: connectex: No connection could be made because the target machine actively refused it.]]
	//test: RegisterEndpoint() -> [status:OK]
	//test: Do() -> [resp:true] [status:Accepted]
	//test: Do() -> [write-requestId:123-request-id] [response-requestId:123-request-id]
	//test: Do() -> [write-relatesTo:test-relates-to] [response-relatesTo:test-relates-to]
	//test: Do() -> [write-statusCode:202] [response-statusCode:202]
	//test: Do() -> [write-content:this is response write content] [response-content:this is response write content]

}

/*
func ExampleDoHttp() {
	req, _ := http.NewRequest(http.MethodGet, "https:/www/google.com/search?q=golang", nil)
	resp, result := DoHttp(req)
	if !result.OK() {
		fmt.Printf("test: DoHttp() -> [loc:%v] [result:%v]\n", result.Location, result)
	} else {
		fmt.Printf("test: DoHttp() -> [status-code:%v] [result:%v]\n", resp.StatusCode, result)
	}

	req, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, result = DoHttp(req)
	if !result.OK() {
		fmt.Printf("test: DoHttp() -> [result:%v]\n", result)
	} else {
		fmt.Printf("test: DoHttp() -> [status-code:%v] [result:%v]\n", resp.StatusCode, result)
	}

	//Output:
}


*/
