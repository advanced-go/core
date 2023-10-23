package httpx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

const (
	ContentTypeJson = "application/json"
	ContentType     = "Content-Type"
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
	WriteResponse[runtime.DebugError](w, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,nil) -> [status:%v] [body:%v]\n", w.Code, w.Body.String())

	w = httptest.NewRecorder()
	WriteResponse[runtime.DebugError](w, []byte(str), nil)
	fmt.Printf("test: WriteResponse(w,%v,nil) -> [status:%v] [body:%v]\n", str, w.Code, w.Body.String())

	//Output:
	//test: WriteResponse(w,nil,nil) -> [status:200] [body:]
	//test: WriteResponse(w,text response,nil) -> [status:200] [body:text response]

}

func ExampleWriteResponse_StatusOK() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatusCode(runtime.StatusOK)
	status.SetMetadata(ContentType, ContentTypeJson)
	WriteResponse[runtime.DebugError](w, []byte(str), status, ContentType)
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [body:%v] [header:%v]\n", str, w.Code, w.Body.String(), resp.Header)

	//Output:
	//test: WriteResponse(w,text response,status) -> [status:200] [body:text response] [header:map[Content-Type:[application/json]]]

}

func ExampleWriteResponse_StatusNotOK() {
	str := "server unavailable"

	w := httptest.NewRecorder()
	status := runtime.NewStatusCode(runtime.StatusUnavailable).SetContent(str)
	WriteResponse[runtime.DebugError](w, nil, status, ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusNotFound).SetContent([]byte("not found"))
	WriteResponse[runtime.DebugError](w, nil, status, ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	str = "operation timed out"
	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusDeadlineExceeded).SetContent(errors.New(str))
	WriteResponse[runtime.DebugError](w, nil, status, ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusInvalidArgument).SetContent(commandTag{
		Sql:          "insert 1",
		RowsAffected: 1,
		Insert:       true,
		Update:       false,
		Delete:       false,
	})
	WriteResponse[runtime.DebugError](w, nil, status, ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	//Output:
	//test: WriteResponse(w,nil,status) -> [status:503] [body:server unavailable] [header:map[Content-Type:[text/plain]]]
	//test: WriteResponse(w,nil,status) -> [status:404] [body:not found] [header:map[Content-Type:[application/json]]]
	//test: WriteResponse(w,nil,status) -> [status:504] [body:operation timed out] [header:map[Content-Type:[text/plain]]]
	//test: WriteResponse(w,nil,status) -> [status:400] [body:{"Sql":"insert 1","RowsAffected":1,"Insert":true,"Update":false,"Delete":false}] [header:map[Content-Type:[application/json]]]

}

func ExampleWriteResponseCopy() {
	w := httptest.NewRecorder()

	resp := http.Response{Header: http.Header{}}
	resp.StatusCode = 504
	resp.Header.Add("key", "value")
	resp.Header.Add("key1", "value1")
	resp.Header.Add("key2", "value2")

	resp.Body = io.NopCloser(bytes.NewReader([]byte("error content")))
	WriteResponseCopy[runtime.DebugError](w, &resp, "key")
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("foo")))
	w = httptest.NewRecorder()
	resp.StatusCode = 200
	WriteResponseCopy[runtime.DebugError](w, &resp, "*")
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	//Output:
	//test: WriteResponse(w,resp,status) -> [status:504] [body:error content] [header:map[Key:[value]]]
	//test: WriteResponse(w,resp,status) -> [status:200] [body:foo] [header:map[Key:[value] Key1:[value1] Key2:[value2]]]

}
