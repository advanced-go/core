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
	s.AddLocation("github.com/advanced-go/location-2")
	s.AddLocation("github.com/advanced-go/location-1")
	if st, ok := any(s).(*status); ok {
		st.errs = append(st.errs, errors.New("test error message 1"), errors.New("testing error msg 2"))
	}
	str := DefaultErrorFormatter(s)
	fmt.Printf("test: DefaultErrorFormatter() -> %v", str)

	//Output:
	//test: DefaultErrorFormatter() -> { "code":200, "status":"OK", "request-id":"1234-5678", "trace" : [ "github.com/advanced-go/location-1","github.com/advanced-go/location-2" ], "errors" : [ "test error message 1","testing error msg 2" ] }

}

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := NewRequestIdContext(nil, "")
	err := errors.New("test error")
	var h LogError

	//s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	s := h.Handle(NewStatus(http.StatusOK), GetOrCreateRequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	//	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	s = h.Handle(NewStatusError(http.StatusInternalServerError, location, err), GetOrCreateRequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = NewStatusError(http.StatusInternalServerError, location)
	fmt.Printf("test: Handle(nil,s) -> [%v] [errors:%v]\n", h.Handle(s, GetOrCreateRequestId(ctx), ""), s.IsErrors())

	s = NewStatusError(http.StatusInternalServerError, location, err)
	errors := s.IsErrors()
	s1 := h.Handle(s, GetOrCreateRequestId(ctx), "")
	fmt.Printf("test: Handle(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.IsErrors())

	//Output:
	//test: Handle(ctx,location,nil) -> [OK] [errors:false]
	//test: Handle(ctx,location,err) -> [Internal Error [test error]] [errors:true]
	//test: Handle(nil,s) -> [OK] [errors:false]
	//test: Handle(nil,s) -> [prev:Internal Error [test error]] [prev-errors:true] [curr:Internal Error [test error]] [curr-errors:true]

}

func Example_InvalidTypeError() {
	fmt.Printf("test: NewInvalidBodyTypeError(nil) -> %v\n", NewInvalidBodyTypeError(nil))
	fmt.Printf("test: NewInvalidBodyTypeError(string) -> %v\n", NewInvalidBodyTypeError("test data"))
	fmt.Printf("test: NewInvalidBodyTypeError(int) -> %v\n", NewInvalidBodyTypeError(500))

	req, _ := http.NewRequest("patch", "https://www.google.com/search", nil)
	fmt.Printf("test: NewInvalidBodyTypeError(*http.Request) -> %v\n", NewInvalidBodyTypeError(req))

	//Output:
	//test: NewInvalidBodyTypeError(nil) -> invalid body type: <nil>
	//test: NewInvalidBodyTypeError(string) -> invalid body type: string
	//test: NewInvalidBodyTypeError(int) -> invalid body type: int
	//test: NewInvalidBodyTypeError(*http.Request) -> invalid body type: *http.Request

}

/*
func ExampleErrorHandleFn() {
	loc := PkgUri + "/ErrorHandleFn"

	fn := NewErrorHandler[LogError]()
	fn("", loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	Output:
	test: Handle[LogErrorHandler]()

}


*/
