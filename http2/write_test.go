package http2

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"net/http/httptest"
)

type data struct {
	Item  string
	Count int
}

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
	WriteResponse[runtime.Output](w, nil, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,nil) -> [status:%v] [body:%v]\n", w.Code, w.Body.String())

	w = httptest.NewRecorder()
	WriteResponse[runtime.Output](w, str, nil, nil)
	fmt.Printf("test: WriteResponse(w,%v,nil) -> [status:%v] [body:%v]\n", str, w.Code, w.Body.String())

	//Output:
	//test: WriteResponse(w,nil,nil) -> [status:200] [body:]
	//test: WriteResponse(w,text response,nil) -> [status:200] [body:text response]

}

func ExampleWriteResponse_StatusOK() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusOK)
	WriteResponse[runtime.Output](w, str, status, nil)
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [body:%v] [header:%v]\n", str, status, w.Body.String(), resp.Header)

	//Output:
	//test: WriteResponse(w,text response,status) -> [status:OK] [body:text response] [header:map[Content-Type:[text/plain; charset=utf-8]]]

}

func ExampleWriteResponse_StatusNotOK() {
	str := "server unavailable"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusServiceUnavailable).SetContent(str, false)
	WriteResponse[runtime.Output](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	w = httptest.NewRecorder()
	status = runtime.NewStatus(http.StatusNotFound).SetContent([]byte("not found"), false)
	WriteResponse[runtime.Output](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	str = "operation timed out"
	w = httptest.NewRecorder()
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded).SetContent(errors.New(str), false)
	WriteResponse[runtime.Output](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	w = httptest.NewRecorder()
	status = runtime.NewStatus(runtime.StatusInvalidArgument).SetContent(commandTag{
		Sql:          "insert 1",
		RowsAffected: 1,
		Insert:       true,
		Update:       false,
		Delete:       false,
	}, false)
	WriteResponse[runtime.Output](w, nil, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	//Output:
	//test: WriteResponse(w,nil,status) -> [status:503] [body:server unavailable] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:404] [body:not found] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:504] [body:operation timed out] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:500] [body:] [header:map[]]

}

func ExampleWriteResponse_Body() {
	w := httptest.NewRecorder()

	body := io.NopCloser(bytes.NewReader([]byte("error content")))
	WriteResponse[runtime.Output](w, body, runtime.NewStatus(http.StatusGatewayTimeout), nil)
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	body = io.NopCloser(bytes.NewReader([]byte("foo")))
	w = httptest.NewRecorder()
	WriteResponse[runtime.Output](w, body, nil,
		[]Attr{{"key", "value"}, {"key1", "value1"}, {"key2", "value2"}})
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	//Output:
	//test: WriteResponse(w,resp,status) -> [status:504] [body:error content] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,resp,status) -> [status:200] [body:foo] [header:map[Content-Type:[text/plain; charset=utf-8] Key:[value] Key1:[value1] Key2:[value2]]]

}

func ExampleWriteStatusContent() {
	r := httptest.NewRecorder()

	// No content
	writeStatusContent[runtime.Output](r, runtime.StatusOK(), "test location")
	r.Result().Header = r.Header()
	buf, status := runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Error message
	r = httptest.NewRecorder()
	writeStatusContent[runtime.Output](r, runtime.NewStatus(http.StatusInternalServerError).SetContent("error message", false), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Json
	d := data{Item: "test item", Count: 500}
	r = httptest.NewRecorder()
	status = runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true)
	//status.SetContent(d, true)

	writeStatusContent[runtime.Output](r, status, "test location") //runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	//Output:
	//test: writeStatusContent() -> 200 [header:map[]] [body:] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Type:[text/plain; charset=utf-8]]] [body:error message] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Type:[application/json]]] [body:{"Item":"test item","Count":500}] [ReadAll:OK]

}
