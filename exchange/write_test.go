package exchange

import (
	"errors"
	"fmt"
	"github.com/go-http-utils/headers"
	"github.com/gotemplates/core/runtime"
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
	WriteResponse(w, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,nil) -> [status:%v] [body:%v]\n", w.Code, w.Body.String())

	w = httptest.NewRecorder()
	WriteResponse(w, []byte(str), nil)
	fmt.Printf("test: WriteResponse(w,%v,nil) -> [status:%v] [body:%v]\n", str, w.Code, w.Body.String())

	//Output:
	//test: WriteResponse(w,nil,nil) -> [status:200] [body:]
	//test: WriteResponse(w,text response,nil) -> [status:200] [body:text response]
	//

}

func ExampleWriteResponse_StatusOK() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatusCode(runtime.StatusOK)
	status.SetMetadata(headers.ContentType, ContentTypeJson)
	WriteResponse(w, []byte(str), status, headers.ContentType)
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [body:%v] [header:%v]\n", str, w.Code, w.Body.String(), resp.Header)

	//Output:
	//test: WriteResponse(w,text response,status) -> [status:200] [body:text response] [header:map[Content-Type:[application/json]]]

}

func ExampleWriteResponse_StatusNotOK() {
	str := "server unavailable"

	w := httptest.NewRecorder()
	status := runtime.NewStatusCode(runtime.StatusUnavailable).SetContent(str)
	WriteResponse(w, nil, status, headers.ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusNotFound).SetContent([]byte("not found"))
	WriteResponse(w, nil, status, headers.ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	str = "operation timed out"
	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusDeadlineExceeded).SetContent(errors.New(str))
	WriteResponse(w, nil, status, headers.ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	w = httptest.NewRecorder()
	status = runtime.NewStatusCode(runtime.StatusInvalidArgument).SetContent(commandTag{
		Sql:          "insert 1",
		RowsAffected: 1,
		Insert:       true,
		Update:       false,
		Delete:       false,
	})
	WriteResponse(w, nil, status, headers.ContentType)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	//Output:
	//test: WriteResponse(w,nil,status) -> [status:503] [body:server unavailable] [header:map[Content-Type:[text/plain]]]
	//test: WriteResponse(w,nil,status) -> [status:404] [body:not found] [header:map[Content-Type:[application/json]]]
	//test: WriteResponse(w,nil,status) -> [status:504] [body:operation timed out] [header:map[Content-Type:[text/plain]]]
	//test: WriteResponse(w,nil,status) -> [status:400] [body:{"Sql":"insert 1","RowsAffected":1,"Insert":true,"Update":false,"Delete":false}] [header:map[Content-Type:[application/json]]]

}
