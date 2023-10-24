package runtime

import (
	"errors"
	"fmt"
)

func ExampleDebugHandler_Handle() {
	location := "/DebugHandler"
	ctx := ContextWithRequestId(nil, "123-request-id")
	err := errors.New("test error")
	var h DebugError

	s := h.Handle(RequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(RequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, nil).SetLocationAndId(location, RequestId(ctx))
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, err).SetLocationAndId(location, RequestId(ctx))
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: Handle(ctx,location,err) -> [Internal] [errors:false]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: HandleStatus(nil,s) -> [prev:Internal] [prev-errors:true] [curr:Internal] [curr-errors:false]

}

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := ContextWithRequestId(nil, "")
	err := errors.New("test error")
	var h LogError

	s := h.Handle(RequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(RequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, nil).SetLocationAndId(location, RequestId(ctx))
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, err).SetLocationAndId(location, RequestId(ctx))
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//test: Handle(ctx,location,err) -> [Internal] [errors:false]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//test: HandleStatus(nil,s) -> [prev:Internal] [prev-errors:true] [curr:Internal] [curr-errors:false]

}

func ExampleErrorHandleFn() {
	loc := PkgUri + "/ErrorHandleFn"
	err := errors.New("debug - error message")

	fn := NewErrorHandler[DebugError]()
	fn("", loc, err)
	fmt.Printf("test: Handle[DebugErrorHandler]()\n")

	fn = NewErrorHandler[LogError]()
	fn("", loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	//Output:
	//[[] github.com/go-ai-agent/core/runtime/ErrorHandleFn [debug - error message]]
	//test: Handle[DebugErrorHandler]()
	//test: Handle[LogErrorHandler]()

}
