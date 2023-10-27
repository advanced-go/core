package exchange

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/exchange/exchangetest"
	"github.com/go-ai-agent/core/httpx"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/runtimetest"
	"io"
	"net/http"
)

const (
	helloWorldUri         = "proxy://www.somestupidname.come"
	serviceUnavailableUri = "http://www.unavailable.com"
	http503FileName       = "file://[cwd]/exchangetest/resource/http-503.txt"
)

// When reading http from a text file, be sure you have the expected blank line between header and body.
// If there is not a blank line after the header section, even if there is not a body, you will receive an
// Unexpected EOF error when calling the golang http.ReadResponse function.
func exchangeProxy(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request or request URL is nil")
	}
	switch exchangetest.Pattern(req) {
	case exchangetest.HttpErrorUri, exchangetest.BodyIOErrorUri:
		return exchangetest.ErrorProxy(req)
	case helloWorldUri:
		resp := exchangetest.NewResponse(http.StatusOK, []byte("<html><body><h1>Hello, World</h1></body></html>"), "content-type", "text/html", "content-length", "1234")
		return resp, nil
	case serviceUnavailableUri:
		// Read the response from an embedded file system.
		//
		// ReadResponseTest(name string)  is only used for calls from do_test.go. When calling from other test
		// files, use the ReadResponse(f fs.FS, name string)
		//
		resp, err := httpx.ReadResponse(httpx.ParseRaw(http503FileName))
		return resp, err
	default:
		fmt.Printf("test: doProxy(req) : unmatched pattern %v", exchangetest.Pattern(req))
	}
	return nil, nil
}

var exchangeCtx = runtime.ContextWithProxy(nil, exchangeProxy)

func ExampleDo_InvalidArgument() {
	_, s := Do[runtimetest.DebugError](nil)
	fmt.Printf("test: Do[runtimetest.DebugError](nil) -> [%v]\n", s)

	//Output:
	//[[] github.com/go-ai-agent/core/exchange/Do [invalid argument : request is nil]]
	//test: Do[runtimetest.DebugError](nil) -> [InvalidArgument github.com/go-ai-agent/core/exchange/Do [invalid argument : request is nil]]

}

func ExampleDo_Proxy_HttpError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, exchangetest.HttpErrorUri, nil)
	resp, err := Do[runtimetest.DebugError](req)
	fmt.Printf("test: Do[runtimetest.DebugError](req) -> [%v] [response:%v]\n", err, resp)

	//Output:
	//[[] github.com/go-ai-agent/core/exchange/Do [http: connection has been hijacked]]
	//test: Do[runtimetest.DebugError](req) -> [Internal Error github.com/go-ai-agent/core/exchange/Do [http: connection has been hijacked]] [response:<nil>]

}

func ExampleDo_Proxy_IOError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, exchangetest.BodyIOErrorUri, nil)
	resp, err := Do[runtimetest.DebugError](req)
	fmt.Printf("test: Do[runtimetest.DebugError](req) -> [%v] [resp:%v] [statusCode:%v] [body:%v]\n", err, resp != nil, resp.StatusCode, resp.Body != nil)

	defer resp.Body.Close()
	buf, s2 := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) -> [%v] [body:%v]\n", s2, string(buf))

	//Output:
	//test: Do[runtimetest.DebugError](req) -> [OK] [resp:true] [statusCode:200] [body:true]
	//test: io.ReadAll(resp.Body) -> [unexpected EOF] [body:]

}

func ExampleDo_Proxy_HellowWorld() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, helloWorldUri, nil)
	resp, err := Do[runtimetest.DebugError](req)
	fmt.Printf("test: Do[runtimetest.DebugError](req) -> [%v] [resp:%v] [statusCode:%v] [content-type:%v] [content-length:%v] [body:%v]\n",
		err, resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Header.Get("content-length"), resp.Body != nil)

	defer resp.Body.Close()
	buf, ioError := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do[runtimetest.DebugError](req) -> [OK] [resp:true] [statusCode:200] [content-type:text/html] [content-length:1234] [body:true]
	//test: io.ReadAll(resp.Body) -> [err:<nil>] [body:<html><body><h1>Hello, World</h1></body></html>]

}

func ExampleDo_Proxy_ServiceUnavailable() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, serviceUnavailableUri, nil)
	resp, _ := Do[runtimetest.DebugError](req)
	fmt.Printf("test: Do[runtimetest.DebugError](req) -> [resp:%v] [statusCode:%v] [content-type:%v] [body:%v]\n",
		resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Body != nil)

	//defer resp.Body.Close()
	//buf, ioError := io.ReadAll(resp.Body)
	//fmt.Printf("test: io.ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do[runtimetest.DebugError](req) -> [resp:true] [statusCode:503] [content-type:text/html] [body:true]

}

func Example_DoT() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := DoT[runtimetest.DebugError, []byte](req)
	fmt.Printf("test: DoT[runtimetest.DebugError,[]byte](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: DoT[runtimetest.DebugError,[]byte](req) -> [status:OK] [buf:true] [resp:true]

}
