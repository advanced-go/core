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

	s = NewStatusError(http.StatusOK, "", errors.New("this is an error message"))
	s.AddLocation("github/advanced-go/core/runtime:AddLocation")
	s.AddLocation("github/advanced-go/core/runtime:TopOfList")
	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error()}, s.Trace(), "1234-56-789"))

	//Output:
	//test: NewStatus() -> [status:OK] [type:github.com/advanced-go/core/runtime/Status]
	//test: NewStatus() -> { "code":200, "status":"OK", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : null }

}

func ExampleNewStatus_Teapot() {
	s := NewStatus(http.StatusTeapot)

	fmt.Printf("test: NewStatus() -> [status:%v]\n", s)

	s = NewStatusError(http.StatusTeapot, "", errors.New("this is an error message"))
	s.AddLocation("github/advanced-go/core/runtime:AddLocation")
	s.AddLocation("github/advanced-go/core/runtime:TopOfList")
	fmt.Printf("test: NewStatus() -> %v\n", defaultFormatter(s.Code, []error{s.Error()}, s.Trace(), "1234-56-789"))

	//Output:
	//test: NewStatus() -> [status:I'm A Teapot]
	//test: NewStatus() -> { "code":418, "status":"I'm A Teapot", "request-id":"1234-56-789", "errors" : [ "this is an error message" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/runtime#TopOfList","https://github.com/advanced-go/core/tree/main/runtime#AddLocation","" ] }

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
