package http2

import (
	"fmt"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/runtime/runtimetest"
	"net/http"
)

func Example_writeStatusContent() {
	r := NewRecorder()

	// No content
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatusOK(), "test location")
	r.Result().Header = r.Header()
	buf, status := io2.ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Error message
	r = NewRecorder()
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatus(http.StatusInternalServerError).SetContent("error message", false), "test location")
	r.Result().Header = r.Header()
	buf, status = io2.ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Json
	d := data{Item: "test item", Count: 500}
	r = NewRecorder()
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true), "test location")
	r.Result().Header = r.Header()
	buf, status = io2.ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	//Output:
	//test: writeStatusContent() -> 200 [header:map[]] [body:] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[13] Content-Type:[text/plain; charset=utf-8]]] [body:error message] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[32] Content-Type:[application/json]]] [body:{"Item":"test item","Count":500}] [ReadAll:OK]

}
