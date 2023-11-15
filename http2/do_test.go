package http2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/http2/http2test"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

const (
	helloWorldUri         = "proxy://www.somestupidname.come"
	serviceUnavailableUri = "http://www.unavailable.com"
	http503FileName       = "file://[cwd]/http2test/resource/http-503.txt"
)

// When reading http from a text file, be sure you have the expected blank line between header and body.
// If there is not a blank line after the header section, even if there is not a body, you will receive an
// Unexpected EOF error when calling the golang http.ReadResponse function.
func exchangeProxy(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request or request URL is nil")
	}
	switch http2test.Pattern(req) {
	case http2test.HttpErrorUri, http2test.BodyIOErrorUri:
		return http2test.ErrorProxy(req)
	case helloWorldUri:
		resp := http2test.NewResponse(http.StatusOK, []byte("<html><body><h1>Hello, World</h1></body></html>"), "content-type", "text/html", "content-length", "1234")
		return resp, nil
	case serviceUnavailableUri:
		// Read the response from an embedded file system.
		//
		// ReadResponseTest(name string)  is only used for calls from do_test.go. When calling from other test
		// files, use the ReadResponse(f fs.FS, name string)
		//
		resp, err := ReadResponse(ParseRaw(http503FileName))
		return resp, err
	default:
		fmt.Printf("test: doProxy(req) : unmatched pattern %v", http2test.Pattern(req))
	}
	return nil, nil
}

var exchangeCtx = runtime.NewProxyContext(nil, exchangeProxy)

func ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [Invalid Argument [invalid argument : request is nil]]

}

func ExampleDo_Proxy_HttpError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, http2test.HttpErrorUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [response:%v]\n", err, resp)

	//Output:
	//test: Do(req) -> [Internal Error [http: connection has been hijacked]] [response:&{internal server error 500  0 0 map[] <nil> 0 [] false false map[] <nil> <nil>}]

}

func ExampleDo_Proxy_IOError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, http2test.BodyIOErrorUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [body:%v]\n", err, resp != nil, resp.StatusCode, resp.Body != nil)

	defer resp.Body.Close()
	buf, s2 := io2.ReadAll(resp.Body)
	fmt.Printf("test: ReadAll(resp.Body) -> [%v] [body:%v]\n", s2, string(buf))

	//Output:
	//test: Do(req) -> [OK] [resp:true] [statusCode:200] [body:true]
	//test: ReadAll(resp.Body) -> [I/O Failure [unexpected EOF]] [body:]

}

func ExampleDo_Proxy_HellowWorld() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, helloWorldUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [content-type:%v] [content-length:%v] [body:%v]\n",
		err, resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Header.Get("content-length"), resp.Body != nil)

	defer resp.Body.Close()
	buf, status := io2.ReadAll(resp.Body)
	fmt.Printf("test: ReadAll(resp.Body) -> [status:%v] [body:%v]\n", status, string(buf))

	//Output:
	//test: Do(req) -> [OK] [resp:true] [statusCode:200] [content-type:text/html] [content-length:1234] [body:true]
	//test: ReadAll(resp.Body) -> [status:OK] [body:<html><body><h1>Hello, World</h1></body></html>]

}

func ExampleDo_Proxy_ServiceUnavailable() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, serviceUnavailableUri, nil)
	resp, _ := Do(req)
	fmt.Printf("test: Do(req) -> [resp:%v] [statusCode:%v] [content-type:%v] [body:%v]\n",
		resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Body != nil)

	//defer resp.Body.Close()
	//buf, ioError := io.ReadAll(resp.Body)
	//fmt.Printf("test: ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do(req) -> [resp:true] [statusCode:503] [content-type:text/html] [body:true]

}

func Example_DoT() {
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=test", nil)
	resp, buf, status := DoT[[]byte](req)
	fmt.Printf("test: DoT[[]byte](req) -> [status:%v] [buf:%v] [resp:%v]\n", status, len(buf) > 0, resp != nil)

	//Output:
	//test: DoT[[]byte](req) -> [status:OK] [buf:true] [resp:true]

}
