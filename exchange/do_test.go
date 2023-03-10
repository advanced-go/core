package exchange

import (
	"context"
	"errors"
	"fmt"
	"github.com/gotemplates/core/exchange/httptest"
	"github.com/gotemplates/core/runtime"
	"io"
	"net/http"
	"time"
)

const (
	helloWorldUri         = "proxy://www.somestupidname.come"
	serviceUnavailableUri = "http://www.unavailable.com"
	http503FileName       = "resource/http/http-503.txt"
)

//var doCtx = ContextWithDo(context.Background(), doProxy)
//type exchangeContext interface {
//	context.Context
//	Exchange
//}

type exchangeData struct {
	ctx  context.Context
	exec HttpExchange
}

func newContext(ctx context.Context, exec HttpExchange) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return &exchangeData{ctx: ctx, exec: exec}
}

func (c *exchangeData) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c *exchangeData) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *exchangeData) Err() error {
	return c.ctx.Err()
}

func (c *exchangeData) Value(key any) any {
	return c.ctx.Value(key)
}

func (c *exchangeData) Do(req *http.Request) (*http.Response, error) {
	if c.exec == nil {
		return nil, errors.New("invalid argument: Exchange interface is nil")
	}
	return c.exec.Do(req)
}

var exchangeCtx = runtime.ContextWithHttpExchange(nil, exchangeProxy)

//
// When reading http from a text file, be sure you have the expected blank line between header and body.
// If there is not a blank line after the header section, even if there is not a body, you will receive an
// Unexpected EOF error when calling the golang http.ReadResponse function.
//

func exchangeProxy(req *http.Request) (*http.Response, error) {
	if req == nil || req.URL == nil {
		return nil, errors.New("request or request URL is nil")
	}
	switch httptest.Pattern(req) {
	case httptest.HttpErrorUri, httptest.BodyIOErrorUri:
		return httptest.ErrorProxy(req)
	case helloWorldUri:
		resp := httptest.NewResponse(http.StatusOK, []byte("<html><body><h1>Hello, World</h1></body></html>"), "content-type", "text/html", "content-length", "1234")
		return resp, nil
	case serviceUnavailableUri:
		// Read the response from an embedded file system.
		//
		// ReadResponseTest(name string)  is only used for calls from do_test.go. When calling from other test
		// files, use the ReadResponse(f fs.FS, name string)
		//
		resp, err := httptest.ReadResponseTest(http503FileName)
		return resp, err
	default:
		fmt.Printf("test: doProxy(req) : unmatched pattern %v", httptest.Pattern(req))
	}
	return nil, nil
}

func ExampleDo_InvalidArgument() {
	_, s := Do(nil)
	fmt.Printf("test: Do(nil) -> [%v]\n", s)

	req, _ := http.NewRequest("", "http://www.google.com", nil)
	_, s = DoClient(req, nil)
	fmt.Printf("test: DoClient(req,nil) -> [%v]\n", s)

	//Output:
	//test: Do(nil) -> [request is nil]
	//test: DoClient(req,nil) -> [client is nil]

}

func ExampleDo_Proxy_HttpError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, httptest.HttpErrorUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [response:%v]\n", err, resp)

	//Output:
	//test: Do(req) -> [http: connection has been hijacked] [response:<nil>]

}

func ExampleDo_Proxy_IOError() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, httptest.BodyIOErrorUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [body:%v]\n", err, resp != nil, resp.StatusCode, resp.Body != nil)

	defer resp.Body.Close()
	buf, s2 := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) -> [%v] [body:%v]\n", s2, string(buf))

	//Output:
	//test: Do(req) -> [<nil>] [resp:true] [statusCode:200] [body:true]
	//test: io.ReadAll(resp.Body) -> [unexpected EOF] [body:]

}

func ExampleDo_Proxy_HellowWorld() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, helloWorldUri, nil)
	resp, err := Do(req)
	fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [content-type:%v] [content-length:%v] [body:%v]\n",
		err, resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Header.Get("content-length"), resp.Body != nil)

	defer resp.Body.Close()
	buf, ioError := io.ReadAll(resp.Body)
	fmt.Printf("test: io.ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do(req) -> [<nil>] [resp:true] [statusCode:200] [content-type:text/html] [content-length:1234] [body:true]
	//test: io.ReadAll(resp.Body) -> [err:<nil>] [body:<html><body><h1>Hello, World</h1></body></html>]

}

func ExampleDo_Proxy_ServiceUnavailable() {
	req, _ := http.NewRequestWithContext(exchangeCtx, http.MethodGet, serviceUnavailableUri, nil)
	resp, err := Do(req)
	if err != nil {
		fmt.Printf("test: Do(req) -> %v\n", err)
	} else {
		fmt.Printf("test: Do(req) -> [%v] [resp:%v] [statusCode:%v] [content-type:%v] [body:%v]\n",
			err, resp != nil, resp.StatusCode, resp.Header.Get("content-type"), resp.Body != nil)
	}
	//defer resp.Body.Close()
	//buf, ioError := io.ReadAll(resp.Body)
	//fmt.Printf("test: io.ReadAll(resp.Body) -> [err:%v] [body:%v]\n", ioError, string(buf))

	//Output:
	//test: Do(req) -> [<nil>] [resp:true] [statusCode:503] [content-type:text/html] [body:true]

}
