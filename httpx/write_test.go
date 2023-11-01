package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/runtimetest"
	"io"
	"net/http"
	"net/http/httptest"
)

type commandTag struct {
	Sql          string
	RowsAffected int64
	Insert       bool
	Update       bool
	Delete       bool
}

func ExampleWriteResponse_NoStatus() {
	str := "text response"

	w := httptest.NewRecorder()
	WriteResponse[runtimetest.DebugError](w, nil, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,nil) -> [status:%v] [body:%v]\n", w.Code, w.Body.String())

	w = httptest.NewRecorder()
	WriteResponse[runtimetest.DebugError](w, str, nil, nil)
	fmt.Printf("test: WriteResponse(w,%v,nil) -> [status:%v] [body:%v]\n", str, w.Code, w.Body.String())

	//Output:
	//test: WriteResponse(w,nil,nil) -> [status:200] [body:]
	//test: WriteResponse(w,text response,nil) -> [status:200] [body:text response]

}

func ExampleWriteResponse_StatusOK() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusOK)
	WriteResponse[runtimetest.DebugError](w, str, status, nil)
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [body:%v] [header:%v]\n", str, status, w.Body.String(), resp.Header)

	//Output:
	//test: WriteResponse(w,text response,status) -> [status:OK] [body:text response] [header:map[]]

}

/*
func ExampleWriteResponse_StatusOK_InvalidKV() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(runtime.StatusOK).SetRequestId("123456-id")
	status1 := WriteResponse[runtimetest.DebugError, string](w, str, status, []Attr{{ContentType,"")
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [status1:%v] [body:%v] [header:%v]\n", str, status1, w.Code, w.Body.String(), resp.Header)

	//Output:
	//{ "id":"123456-id", "l":"github.com/go-ai-agent/core/httpx/WriteResponse", "o":null "err" : [ "invalid number of kv items: number is odd, missing a value" ] }
	//test: WriteResponse(w,text response,status) -> [status:Internal [invalid number of kv items: number is odd, missing a value]] [status1:500] [body:] [header:map[]]

}


*/

func ExampleWriteResponse_StatusNotOK() {
	str := "server unavailable"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusServiceUnavailable).SetContent(str, false)
	WriteResponse[runtimetest.DebugError](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	w = httptest.NewRecorder()
	status = runtime.NewStatus(http.StatusNotFound).SetContent([]byte("not found"), false)
	WriteResponse[runtimetest.DebugError](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	str = "operation timed out"
	w = httptest.NewRecorder()
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded).SetContent(errors.New(str), false)
	WriteResponse[runtimetest.DebugError](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	w = httptest.NewRecorder()
	status = runtime.NewStatus(runtime.StatusInvalidArgument).SetContent(commandTag{
		Sql:          "insert 1",
		RowsAffected: 1,
		Insert:       true,
		Update:       false,
		Delete:       false,
	}, false)
	WriteResponse[runtimetest.DebugError](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	//Output:
	//test: WriteResponse(w,nil,status) -> [status:503] [body:server unavailable] [header:map[Content-Length:[18] Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:404] [body:not found] [header:map[Content-Length:[9] Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:504] [body:operation timed out] [header:map[Content-Length:[19] Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:500] [body:] [header:map[]]

}

func Example_RequestBody() {
	w := httptest.NewRecorder()

	body := io.NopCloser(bytes.NewReader([]byte("error content")))
	WriteResponse[runtimetest.DebugError](w, body, runtime.NewStatus(http.StatusGatewayTimeout), nil)
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	body = io.NopCloser(bytes.NewReader([]byte("foo")))
	w = httptest.NewRecorder()
	WriteResponse[runtimetest.DebugError](w, body, nil,
		[]Attr{{"key", "value"}, {"key1", "value1"}, {"key2", "value2"}})
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	//Output:
	//test: WriteResponse(w,resp,status) -> [status:504] [body:error content] [header:map[Content-Length:[13] Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,resp,status) -> [status:200] [body:foo] [header:map[Content-Length:[3] Content-Type:[text/plain; charset=utf-8] Key:[value] Key1:[value1] Key2:[value2]]]

}
