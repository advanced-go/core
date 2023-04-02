package runtime

import (
	"errors"
	"fmt"
)

func ExampleNoOpErrorHandler3_Handle() {
	location := "/test"
	err := errors.New("test error")
	var h NoOpError3

	fmt.Printf("test: Handle(location,nil) -> [%v]\n", h.Handle(location, nil))
	fmt.Printf("test: Handle(location,err) -> [%v]\n", h.Handle(location, err))

	s := NewStatus(StatusInternal, location, nil)
	fmt.Printf("test: HandleStatus(s) -> [%v]\n", h.HandleStatus(s))

	s = NewStatus(StatusInternal, location, err)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [curr:%v]\n", s, h.HandleStatus(s))

	//Output:
	//test: Handle(location,nil) -> [OK]
	//test: Handle(location,err) -> [Internal [test error]]
	//test: HandleStatus(s) -> [OK]
	//test: HandleStatus(s) -> [prev:Internal [test error]] [curr:Internal [test error]]

}

func ExampleDebugHandler3_Handle() {
	location := "/DebugHandler"
	ctx := ContextWithRequestId(nil, "123-request-id")
	err := errors.New("test error")
	var h DebugError3

	s := h.Handle(location, nil)
	fmt.Printf("test: Handle(location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.HandleWithContext(ctx, location, err)
	fmt.Printf("test: HandleWithContext(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, location, nil).SetContext(ctx)
	fmt.Printf("test: HandleStatus(s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, location, err).SetContext(ctx)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(location,nil) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: HandleWithContext(ctx,location,err) -> [Internal] [errors:false]
	//test: HandleStatus(s) -> [OK] [errors:false]
	//[123-request-id /DebugHandler [test error]]
	//test: HandleStatus(s) -> [prev:Internal] [prev-errors:true] [curr:Internal] [curr-errors:false]

}

func ExampleLogHandler3_Handle() {
	location := "/LogHandler"
	ctx := ContextWithRequestId(nil, "")
	err := errors.New("test error")
	var h LogError3

	s := h.Handle(location, nil)
	fmt.Printf("test: Handle(location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.HandleWithContext(ctx, location, err)
	fmt.Printf("test: HandleWithContext(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatus(StatusInternal, location, nil).SetContext(ctx)
	fmt.Printf("test: HandleStatus(s) -> [%v] [errors:%v]\n", h.HandleStatus(s), s.IsErrors())

	s = NewStatus(StatusInternal, location, err).SetContext(ctx)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s)
	fmt.Printf("test: HandleStatus(s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(location,nil) -> [OK] [errors:false]
	//test: HandleWithContext(ctx,location,err) -> [Internal] [errors:false]
	//test: HandleStatus(s) -> [OK] [errors:false]
	//test: HandleStatus(s) -> [prev:Internal] [prev-errors:true] [curr:Internal] [curr-errors:false]

}

func ExampleErrorHandleFn() {
	loc := PkgUrl + "/ErrorHandleFn"
	err := errors.New("debug - error message")

	fn := Handle[DebugError3]()
	fn(loc, err)
	fmt.Printf("test: Handle[DebugErrorHandler]()\n")

	fn = Handle[LogError3]()
	fn(loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	//Output:
	//[[] github.com/go-sre/core/runtime/ErrorHandleFn [debug - error message]]
	//test: Handle[DebugErrorHandler]()
	//test: Handle[LogErrorHandler]()

}

func ExampleErrorHandleStatus() {
	loc := PkgUrl + "/ErrorHandleStatus"
	err := errors.New("debug - error message")

	fn := StatusHandle[DebugError3]()
	status := fn(NewStatusError(loc, err))
	fmt.Printf("test: StatusHandle[DebugErrorHandler]() [status:%v] [errors:%v]\n", status, status.IsErrors())

	fn = StatusHandle[LogError3]()
	status = fn(NewStatusError(loc, errors.New("log - error message")))
	fmt.Printf("test: StatusHandle[LogErrorHandler]() [status:%v] [errors:%v]\n", status, status.IsErrors())

	//Output:
	//[[] github.com/go-sre/core/runtime/ErrorHandleStatus [debug - error message]]
	//test: StatusHandle[DebugErrorHandler]() [status:Internal] [errors:false]
	//test: StatusHandle[LogErrorHandler]() [status:Internal] [errors:false]

}
