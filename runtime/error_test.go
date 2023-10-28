package runtime

import (
	"errors"
	"fmt"
)

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := ContextWithRequestId(nil, "")
	err := errors.New("test error")
	var h LogError

	s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatusError(StatusInternal, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s, GetOrCreateRequestId(ctx), ""), s.IsErrors())

	s = NewStatusError(StatusInternal, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s, GetOrCreateRequestId(ctx), "")
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//test: Handle(ctx,location,err) -> [Internal [test error]] [errors:true]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//test: HandleStatus(nil,s) -> [prev:Internal [test error]] [prev-errors:true] [curr:Internal [test error]] [curr-errors:true]

}

/*
func ExampleErrorHandleFn() {
	loc := PkgUri + "/ErrorHandleFn"

	fn := NewErrorHandler[LogError]()
	fn("", loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	//Output:
	//test: Handle[LogErrorHandler]()

}


*/
