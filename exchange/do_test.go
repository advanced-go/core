package exchange

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"time"
)

const (
	testContent = "this is response write content"
	requestId   = "123-request-id"
	relatesTo   = "test-relates-to"
	statusCode  = http.StatusAccepted
)

// dial tcp [2607:f8b0:4023:1009::68]:443: i/o timeout]

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
		resp != nil, status.Code, status.Error(), resp.Header.Get("content-type"), resp.Body != nil)

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

/*
func _ExampleDo_Proxy() {
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

	buf, _ := io.ReadAll(resp.Body)
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


*/

func defaultDo(r *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		if DeadlineExceededError(r) {
			return &http.Response{StatusCode: http.StatusGatewayTimeout}, err
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, err
	}
	//time.Sleep(time.Second * 2)
	buf, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		if DeadlineExceededError(err1) {
			return &http.Response{StatusCode: http.StatusGatewayTimeout}, err1
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, err1
	}
	if buf != nil {
	}
	return resp, nil
}

func ExampleDefaultDo_OK() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := defaultDo(req)
	fmt.Printf("test: DefaultDo_OK()-no-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*4)
	defer cancel()
	r2 := req.Clone(ctx)
	resp, err = defaultDo(r2)
	fmt.Printf("test: DefaultDo_OK()-5s-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: DefaultDo_OK()-no-timeout -> [status-code:200] [err:<nil>]
	//test: DefaultDo_OK()-5s-timeout -> [status-code:200] [err:<nil>]

}

func ExampleDefaultDo_Timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := defaultDo(req)
	fmt.Printf("test: DefaultDo_Timeout()-Get()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Millisecond*600)
	defer cancel2()
	req, _ = http.NewRequestWithContext(ctx2, http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, err = defaultDo(req)
	fmt.Printf("test: DefaultDo_Timeout()-ReadAll()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: DefaultDo_Timeout()-Get()-timeout -> [status-code:504] [err:Get "https://www.google.com/search?q=golang": context deadline exceeded]
	//test: DefaultDo_Timeout()-ReadAll()-timeout -> [status-code:504] [err:context deadline exceeded]

}

func exchangeDo(r *http.Request) (*http.Response, *runtime.Status) {
	resp, status := Do(r)
	if !status.OK() {
		return resp, status
	}
	buf, err1 := io.ReadAll(resp.Body)
	if err1 != nil {
		if DeadlineExceededError(err1) {
			return &http.Response{StatusCode: http.StatusGatewayTimeout}, runtime.NewStatusError(http.StatusGatewayTimeout, err1)
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, runtime.NewStatusError(http.StatusInternalServerError, err1)
	}
	if buf != nil {
	}
	return resp, nil
}

func ExampleExchangeDo_OK() {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := exchangeDo(req)
	fmt.Printf("test: ExchangeDo_OK()-no-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*4)
	defer cancel()
	r2 := req.Clone(ctx)
	resp, err = exchangeDo(r2)
	fmt.Printf("test: ExchangeDo_OK()-5s-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: ExchangeDo_OK()-no-timeout -> [status-code:200] [err:<nil>]
	//test: ExchangeDo_OK()-5s-timeout -> [status-code:200] [err:<nil>]

}

func ExampleExchangeDo_Timeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.google.com/search?q=golang", nil)

	resp, err := exchangeDo(req)
	fmt.Printf("test: ExchangeDo_Timeout()-Get()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Millisecond*600)
	defer cancel2()
	req, _ = http.NewRequestWithContext(ctx2, http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, err = exchangeDo(req)
	fmt.Printf("test: ExchangeDo_Timeout()-ReadAll()-timeout -> [status-code:%v] [err:%v]\n", resp.StatusCode, err)

	//Output:
	//test: ExchangeDo_Timeout()-Get()-timeout -> [status-code:504] [err:Get "https://www.google.com/search?q=golang": context deadline exceeded]
	//test: ExchangeDo_Timeout()-ReadAll()-timeout -> [status-code:504] [err:Timeout [context deadline exceeded]]

}
