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

	s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatusError(StatusInternal, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s, GetOrCreateRequestId(ctx)), s.IsErrors())

	s = NewStatusError(StatusInternal, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s, GetOrCreateRequestId(ctx))
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: Handle(ctx,location,err) -> [Internal /DebugHandler [test error]] [errors:true]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: HandleStatus(nil,s) -> [prev:Internal /DebugHandler [test error]] [prev-errors:true] [curr:Internal /DebugHandler [test error]] [curr-errors:true]

}

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := ContextWithRequestId(nil, "")
	err := errors.New("test error")
	var h LogError

	s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatusError(StatusInternal, location) //.SetLocationAndId(location, )
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s, GetOrCreateRequestId(ctx)), s.IsErrors())

	s = NewStatusError(StatusInternal, location, err) //.SetLocationAndId(location, GetOrCreateRequestId(ctx))
	errors := s.IsErrors()
	s1 := h.HandleStatus(s, GetOrCreateRequestId(ctx))
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//test: Handle(ctx,location,err) -> [Internal /LogHandler [test error]] [errors:true]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//test: HandleStatus(nil,s) -> [prev:Internal /LogHandler [test error]] [prev-errors:true] [curr:Internal /LogHandler [test error]] [curr-errors:true]

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
