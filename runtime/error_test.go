package runtime

import (
	"errors"
	"fmt"
	"net/http"
)

func Example_FormatUri() {
	s := "github/advanced-go/core/runtime:testFunc"

	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	s = "gitlab/advanced-go/core/runtime:testFunc"
	fmt.Printf("test: formatUri(%v) -> %v\n", s, formatUri(s))

	//Output:
	//test: formatUri(github/advanced-go/core/runtime:testFunc) -> https://github.com/advanced-go/core/tree/main/runtime#testFunc
	//test: formatUri(gitlab/advanced-go/core/runtime:testFunc) -> gitlab/advanced-go/core/runtime:testFunc

}

func Example_FormatUri_Test() {
	s := "http://localhost:8080/github.com/advanced-go/core/runtime/testFunc"
	req, err := http.NewRequest("", s, nil)
	fmt.Printf("test: http.URL -> [req:%v] [url:%v] [err:%v]\n", req != nil, req.URL, err)

	s = "http://localhost:8080/github.com/advanced-go/core/runtime:testFunc"
	req, err = http.NewRequest("", s, nil)
	fmt.Printf("test: http.URL -> [req:%v] [url:%v] [err:%v]\n", req != nil, req.URL, err)

	s = "http://localhost:8080/github.com:advanced-go/core/runtime.testFunc"
	req, err = http.NewRequest("", s, nil)
	fmt.Printf("test: http.URL -> [req:%v] [url:%v] [err:%v]\n", req != nil, req.URL, err)

	//Output:
	//test: http.URL -> [req:true] [url:http://localhost:8080/github.com/advanced-go/core/runtime/testFunc] [err:<nil>]
	//test: http.URL -> [req:true] [url:http://localhost:8080/github.com/advanced-go/core/runtime:testFunc] [err:<nil>]
	//test: http.URL -> [req:true] [url:http://localhost:8080/github.com:advanced-go/core/runtime.testFunc] [err:<nil>]

}

func Example_DefaultFormat() {
	s := NewStatusError(http.StatusNotFound, errors.New("test error message 1"), nil)

	// Adding in reverse to mirror call stack
	//s.AddLocation("github/advanced-go/location-2")
	//s.AddLocation("github/advanced-go/location-1")

	//if st, ok := any(s).(*statusState); ok {
	//	st.Errs = append(st.Errs, errors.New("test error message 1"), errors.New("testing error msg 2"))
	//}
	//SetOutputFormatter()
	str := formatter(s.Code, []error{s.Error()}, s.Trace(), s.Content(), "1234-5678")
	fmt.Printf("test: formatter() -> %v", str)

	//Output:
	//test: formatter() -> { "code":404, "status":"Not Found", "request-id":"1234-5678", "errors" : [ "test error message 1" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#Example_DefaultFormat" ] }

}

func ExampleOutputHandler_Handle() {
	//location := "/OutputHandler"
	//origin := "github.com/module/package/calling-fn"
	ctx := NewRequestIdContext(nil, "123-request-id")
	err := errors.New("test error")
	var h Output

	//status := runtime.NewStatusError(0, location, err)
	s := h.Handle(NewStatus(http.StatusInternalServerError), RequestId(ctx))
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.Error() != nil)

	s = h.Handle(NewStatusError(http.StatusInternalServerError, err, nil), GetOrCreateRequestId2(ctx))
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [handled:%v]\n", s, s.handled)

	s = NewStatusError(http.StatusInternalServerError, nil, nil)
	fmt.Printf("test: HandleStatus(nil,s) -> [%v] [handled:%v]\n", h.Handle(nil, GetOrCreateRequestId2(ctx)), s.handled)

	//Output:
	//test: Handle(ctx,location,nil) -> [Internal Error] [errors:false]
	//{ "code":500, "status":"Internal Error", "request-id":"123-request-id", "errors" : [ "test error" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#ExampleOutputHandler_Handle" ] }
	//test: Handle(ctx,location,err) -> [Internal Error [test error]] [handled:true]
	//test: HandleStatus(nil,s) -> [OK] [handled:false]

}

func ExampleLogHandler_Handle() {
	//location := "/LogHandler"
	ctx := NewRequestIdContext(nil, "")
	err := errors.New("test error")
	var h Log

	//s := h.Handle(GetOrCreateRequestId(ctx), location, nil)
	s := h.Handle(NewStatus(http.StatusOK), GetOrCreateRequestId2(ctx))
	fmt.Printf("test: Handle(ctx,location,nil) -> [%v] [errors:%v]\n", s, s.Error() != nil)

	//	s = h.Handle(GetOrCreateRequestId(ctx), location, err)
	s = h.Handle(NewStatusError(http.StatusInternalServerError, err, nil), GetOrCreateRequestId2(ctx))
	fmt.Printf("test: Handle(ctx,location,err) -> [%v] [errors:%v]\n", s, s.Error() != nil)

	s = NewStatusError(http.StatusInternalServerError, nil, nil)
	fmt.Printf("test: Handle(nil,s) -> [%v] [errors:%v]\n", h.Handle(nil, GetOrCreateRequestId2(ctx)), s.Error() != nil)

	s = NewStatusError(http.StatusInternalServerError, err, nil)
	errors := s.Error() != nil
	s1 := h.Handle(s, GetOrCreateRequestId2(ctx))
	fmt.Printf("test: Handle(nil,s) -> [prev:%v] [prev-errors:%v] [curr:%v] [curr-errors:%v]\n", s, errors, s1, s1.Error() != nil)

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

// ErrorHandleFn - function type for error handling
//type ErrorHandleFn func(requestId, location string, errs ...error) *Status
// NewErrorHandler - templated function providing an error handle function via a closure

func NewErrorHandler[E ErrorHandler]() ErrorHandleFn {
	var e E
	return func(requestId string, location string, errs ...error) *Status {
		return e.Handle(NewStatusError(http.StatusInternalServerError, location, errs...), requestId, "")
	}
}

func ExampleErrorHandleFn() {
	loc := PkgUri + "/ErrorHandleFn"

	fn := NewErrorHandler[LogError]()
	fn("", loc, errors.New("log - error message"))
	fmt.Printf("test: Handle[LogErrorHandler]()\n")

	Output:
	test: Handle[LogErrorHandler]()

}


*/
