package runtime

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

func ExampleNewStatus_OK() {
	s := StatusOK()

	path := reflect.TypeOf(Status{}).PkgPath()
	path += "/" + reflect.TypeOf(Status{}).Name()
	fmt.Printf("test: NewStatus() -> [status:%v] [type:%v]\n", s, path)

	s = NewStatusError(http.StatusBadGateway, errors.New("this is an error message"), nil)
	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error()}, s.Trace(), s.Content(), "1234-56-789"))

	//Output:
	//test: NewStatus() -> [status:OK] [type:github.com/advanced-go/core/runtime/Status]
	//test: NewStatus() -> { "code":502, "status":"error: code not mapped: 502", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#ExampleNewStatus_OK" ] }

}

func ExampleNewStatus_Teapot() {
	s := NewStatus(http.StatusTeapot)

	fmt.Printf("test: NewStatus() -> [status:%v]\n", s)

	s = NewStatusError(http.StatusTeapot, errors.New("this is an error message"), nil)
	//s.AddLocation()
	//s.AddLocation("github/advanced-go/core/runtime:TopOfList")
	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error()}, s.Trace(), s.Content(), "1234-56-789"))

	//Output:
	//test: NewStatus() -> [status:I'm A Teapot]
	//test: NewStatus() -> { "code":418, "status":"I'm A Teapot", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#ExampleNewStatus_Teapot" ] }

}

func ExampleNewStatus_Location() {
	s := errorFunc()
	s.AddLocation()

	str := formatter(s.Code, []error{s.Error()}, s.Trace(), s.Content(), "1234-5678")
	fmt.Printf("test: Location() -> [out:%v] [trace:%v]\n", str, s.Trace())

	//Output:
	//test: Location() -> [out:{ "code":400, "status":"Bad Request", "request-id":"1234-5678", "errors" : [ "test bad request error" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#ExampleNewStatus_Location","https://github.com/advanced-go/core/tree/main/runtime#errorFunc" ] }
	//] [trace:[github/advanced-go/core/runtime:errorFunc github/advanced-go/core/runtime:ExampleNewStatus_Location]]

}

func errorFunc() *Status {
	return NewStatusError(http.StatusBadRequest, errors.New("test bad request error"), nil)
}

/*
func ExampleNewStatus_TeapotHandled() {
	var e Output
	s := NewStatus(http.StatusTeapot)

	//fmt.Printf("test: NewStatus() -> [status:%v]\n", s)

	s.Error = errors.New("this is an error message")
	s.AddLocation("github/advanced-go/core/runtime:AddLocation")
	s.AddLocation("github/advanced-go/core/runtime:TopOfList")

	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error}, s.Trace(), "1234-56-789"))
    //e.Handle()
	//Output:
	//test: NewStatus() -> [status:I'm A Teapot]
	//test: NewStatus() -> { "code":418, "status":"I'm A Teapot", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#TopOfList","https://github.com/advanced-go/core/tree/main/runtime#AddLocation" ] }

}


*/
