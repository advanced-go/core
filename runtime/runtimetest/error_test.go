package runtimetest

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

func Example_formatErrors() {
	fmt.Printf("test: formatErrors(nil) -> %v\n", runtime.FormatErrors("err", nil))

	//Output:
	//test: formatErrors(nil) -> "err" : null

}

func ExampleDebugHandler_Handle() {
	location := "/DebugHandler"
	origin := "github.com/module/package/calling-fn"
	ctx := runtime.ContextWithRequestId(nil, "123-request-id")
	err := errors.New("test error")
	var h DebugError

	//status := runtime.NewStatusError(0, location, err)
	s := h.Handle(runtime.ContextRequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(runtime.GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [handled:%v]\n", s, s.ErrorsHandled())

	s = runtime.NewStatusError(runtime.StatusInternal, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [handled:%v]\n", h.HandleStatus(s, runtime.GetOrCreateRequestId(ctx), origin), s.ErrorsHandled())

	//s = runtime.NewStatusError(runtime.StatusInternal, location, err)
	//errors := s.IsErrors()
	//s1 := h.HandleStatus(s, runtime.GetOrCreateRequestId(ctx), "")
	//fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//{ "id":"123-request-id", "l":"/DebugHandler", "o":null "err" : [ "test error" ] }
	//test: Handle(ctx,location,err) -> [Internal [test error]] [handled:true]
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
