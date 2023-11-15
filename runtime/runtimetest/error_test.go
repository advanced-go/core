package runtimetest

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func Example_formatErrors() {
	s := runtime.NewStatusOK()
	s.SetRequestId("1234-5678")
	// Adding on reverse to mirror call stack
	s.AddLocation("github.com/advanced-go/location-2")
	s.AddLocation("github.com/advanced-go/location-1")
	fmt.Printf("test: defaultFormatter() -> %v", defaultFormatter(s))

	//Output:
	//test: defaultFormatter() -> { "code":200, "status":"OK", "id":"1234-5678", "trace" : [ "github.com/advanced-go/location-1","github.com/advanced-go/location-2" ], "err" : null }

}

func ExampleDebugHandler_Handle() {
	location := "/DebugHandler"
	origin := "github.com/module/package/calling-fn"
	ctx := runtime.NewRequestIdContext(nil, "123-request-id")
	err := errors.New("test error")
	var h DebugError

	//status := runtime.NewStatusError(0, location, err)
	s := h.Handle(runtime.NewStatus(http.StatusInternalServerError), runtime.RequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(runtime.NewStatusError(http.StatusInternalServerError, location, err), runtime.GetOrCreateRequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [handled:%v]\n", s, s.ErrorsHandled())

	s = runtime.NewStatusError(http.StatusInternalServerError, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [handled:%v]\n", h.Handle(s, runtime.GetOrCreateRequestId(ctx), origin), s.ErrorsHandled())

	//s = runtime.NewStatusError(runtime.StatusInternal, location, err)
	//errors := s.IsErrors()
	//s1 := h.HandleStatus(s, runtime.GetOrCreateRequestId(ctx), "")
	//fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [Internal Error] [errors:false]
	//{ "code":500, "status":"Internal Error", "id":"123-request-id", "trace" : [ "/DebugHandler","/DebugHandler" ], "err" : [ "test error" ] }
	//test: Handle(ctx,location,err) -> [Internal Error [test error]] [handled:true]
	//test: HandleStatus(nil,s) -> [OK] [handled:false]

}

/*
func ExampleErrorHandleFn() {
	loc := "/ErrorHandleFn"
	err := errors.New("debug - error message")

	fn := runtime.NewErrorHandler[DebugError]()
	fn("", loc, err)
	fmt.Printf("test: Handle[DebugError]()\n")

	//Output:
	//{ "id":null, "l":"/ErrorHandleFn", "o":null "err" : [ "debug - error message" ] }
	//test: Handle[DebugError]()

}


*/
