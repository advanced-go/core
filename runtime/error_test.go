package runtime

import (
	"errors"
	"fmt"
	"net/http"
)

func Example_FormatUri() {
	s := "github.com/advanced-go/core/runtime:testFunc"

	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	s = "gitlab.com/advanced-go/core/runtime:testFunc"
	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	//Output:
	//test: formatUri(github.com/advanced-go/core/runtime:testFunc) -> github.com/advanced-go/core/tree/main/runtime:testFunc
	//test: formatUri(gitlab.com/advanced-go/core/runtime:testFunc) -> gitlab.com/advanced-go/core/runtime:testFunc

}

func Example_FormatUri_Test() {
	s := "github.com/advanced-go/core/runtime/testFunc"

	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	s = "gitlab.com/advanced-go/core/runtime:testFunc"
	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	//Output:
	//test: formatUri(github.com/advanced-go/core/runtime:testFunc) -> github.com/advanced-go/core/tree/main/runtime:testFunc
	//test: formatUri(gitlab.com/advanced-go/core/runtime:testFunc) -> gitlab.com/advanced-go/core/runtime:testFunc

}

func Example_DefaultFormat() {
	s := NewStatus(http.StatusNotFound)
	s.SetRequestId("1234-5678")
	// Adding in reverse to mirror call stack
	s.AddLocation("github.com/advanced-go/location-2")
	s.AddLocation("github.com/advanced-go/location-1")
	if st, ok := any(s).(*statusState); ok {
		st.Errs = append(st.Errs, errors.New("test error message 1"), errors.New("testing error msg 2"))
	}
	//SetOutputFormatter()
	str := formatter(s)
	fmt.Printf("test: formatter() -> %v", str)

	//Output:
	//test: formatter() -> { "code":404, "status":"Not Found", "request-id":"1234-5678", "trace" : [ "github.com/advanced-go/tree/main/location-1","github.com/advanced-go/tree/main/location-2" ], "errors" : [ "test error message 1","testing error msg 2" ] }

}

func ExampleOutputHandler_Handle() {
	location := "/OutputHandler"
	origin := "github.com/module/package/calling-fn"
	ctx := NewRequestIdContext(nil, "123-request-id")
	err := errors.New("test error")
	var h Output

	//status := runtime.NewStatusError(0, location, err)
	s := h.Handle(NewStatus(http.StatusInternalServerError), RequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.IsErrors())

	s = h.Handle(NewStatusError(http.StatusInternalServerError, location, err), GetOrCreateRequestId(ctx), location)
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [handled:%v]\n", s, errorsHandled(s))

	s = NewStatusError(http.StatusInternalServerError, location)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [handled:%v]\n", h.Handle(s, GetOrCreateRequestId(ctx), origin), errorsHandled(s))

	//Output:
	//test: Handle(ctx,location,nil) -> [Internal Error] [errors:false]
	//{ "code":500, "status":"Internal Error", "request-id":"123-request-id", "trace" : [ "/OutputHandler","/OutputHandler" ], "errors" : [ "test error" ] }
	//test: Handle(ctx,location,err) -> [Internal Error [test error]] [handled:true]
	//test: HandleStatus(nil,s) -> [OK] [handled:false]

}

func ExampleLogHandler_Handle() {
	location := "/LogHandler"
	ctx := NewRequestIdContext(nil, "")
	err := errors.New("test error")
	var h Log

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
