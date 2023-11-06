package runtime

import (
	"errors"
	"fmt"
	"net/http"
)

func Example_DefaultErrorFormat() {
	s := NewStatusOK()
	s.SetRequestId("1234-5678")
	// Adding on reverse to mirror call stack
	s.AddLocation("github.com/go-ai-agent/location-2")
	s.AddLocation("github.com/go-ai-agent/location-1")
	s.errs = append(s.errs, errors.New("test error message 1"), errors.New("testing error msg 2"))
	str := DefaultErrorFormatter(s)
	fmt.Printf("test: DefaultErrorFormatter() -> %v", str)

	//Output:
	//test: DefaultErrorFormatter() -> { "Code":200, "Status":"OK", "RequestId":"1234-5678", "Location" : [ "github.com/go-ai-agent/location-1","github.com/go-ai-agent/location-2" ], "Errors" : [ "test error message 1","testing error msg 2" ] }

}

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := NewRequestIdContext(nil, "")
	err := errors.New("test error")
	var h LogError

	s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatusError(http.StatusInternalServerError, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [errors:%v]\n", h.HandleStatus(s, GetOrCreateRequestId(ctx), ""), s.IsErrors())

	s = NewStatusError(http.StatusInternalServerError, location, err)
	errors := s.IsErrors()
	s1 := h.HandleStatus(s, GetOrCreateRequestId(ctx), "")
	fmt.Printf("test: HandleStatus(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//test: Handle(ctx,location,err) -> [Internal Error [test error]] [errors:true]
	//test: HandleStatus(nil,s) -> [OK] [errors:false]
	//test: HandleStatus(nil,s) -> [prev:Internal Error [test error]] [prev-errors:true] [curr:Internal Error [test error]] [curr-errors:true]

}

func ExampleErrorHandleFn() {
	loc := PkgUri + "/ErrorHandleFn"

	fn := NewErrorHandler[LogError]()
	fn("", loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	//Output:
	//test: Handle[LogErrorHandler]()

}
