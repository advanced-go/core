package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/http/httptest"
)

func Example_writeStatusContent() {
	r := httptest.NewRecorder()

	// No content
	writeStatusContent[runtime.Output](r, runtime.StatusOK(), "test location")
	r.Result().Header = r.Header()
	buf, status := runtime.NewBytes(r.Result())
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [NewBytes:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Error message
	r = httptest.NewRecorder()
	writeStatusContent[runtime.Output](r, runtime.NewStatus(http.StatusInternalServerError).SetContent("error message", false), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.NewBytes(r.Result())
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [NewBytes:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Json
	d := data{Item: "test item", Count: 500}
	r = httptest.NewRecorder()
	writeStatusContent[runtime.Output](r, runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.NewBytes(r.Result())
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [NewBytes:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	//Output:
	//test: writeStatusContent() -> 200 [header:map[]] [body:] [NewBytes:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[13] Content-Type:[text/plain; charset=utf-8]]] [body:error message] [NewBytes:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[32] Content-Type:[application/json]]] [body:{"Item":"test item","Count":500}] [NewBytes:OK]

}
