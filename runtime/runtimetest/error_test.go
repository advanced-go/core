package runtimetest

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

func ExampleDebugHandler_Handle() {
	location := "/DebugHandler"
	ctx := runtime.ContextWithRequestId(nil, "123-request-id")
	err := errors.New("test error")
	var h DebugError

	s := h.Handle(runtime.GetOrCreateRequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(runtime.GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = runtime.NewStatusError(runtime.StatusInternal, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s, runtime.GetOrCreateRequestId(ctx)), s.IsErrors())

	s = runtime.NewStatusError(runtime.StatusInternal, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s, runtime.GetOrCreateRequestId(ctx))
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: Handle(ctx,location,err) -> [Internal /DebugHandler [test error]] [errors:true]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: HandleStatus(nil,s) -> [prev:Internal /DebugHandler [test error]] [prev-errors:true] [curr:Internal /DebugHandler [test error]] [curr-errors:true]

}

func ExampleErrorHandleFn() {
	loc := "/ErrorHandleFn"
	err := errors.New("debug - error message")

	fn := runtime.NewErrorHandler[DebugError]()
	fn("", loc, err)
	fmt.Printf("test: Handle[DebugErrorHandler]()\n")

	//Output:
	//[[] /ErrorHandleFn [debug - error message]]
	//test: Handle[DebugErrorHandler]()

}
