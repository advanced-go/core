package runtime

import (
	"errors"
	"fmt"
	"net/http"
)

type address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (a address) GetZip() string {
	return a.Zip
}

/*
	func ExampleStatus_GRPC_String() {
		s := NewStatus(StatusPermissionDenied)
		fmt.Printf("test: NewStatus() -> [%v]\n", s)

		s = NewStatusError(StatusOutOfRange, "", errors.New("error - 1"), errors.New("error - 2"))
		fmt.Printf("test: NewStatus() -> [%v]\n", s)

		//Output:
		//test: NewStatus() -> [PermissionDenied]
		//test: NewStatus() -> [OutOfRange [error - 1 error - 2]]

}
*/
func Example_NewStatusError() {
	location := "test"
	err := errors.New("http error")
	fmt.Printf("test: NewStatusError(nil) -> [%v]\n", NewStatusError(StatusInvalidContent, location, err))
	fmt.Printf("test: NewStatusError(nil) -> [%v]\n", NewStatusError(http.StatusBadRequest, location, err))

	//resp := http.Response{StatusCode: http.StatusBadRequest}
	//fmt.Printf("test: NewHttpStatus(resp) -> [%v]\n", NewHttpStatus(&resp, nil).SetLocation(location))
	//fmt.Printf("test: NewHttpStatus(resp) -> [%v]\n", NewHttpStatus(&resp, err).SetLocation(location))

	//Output:
	//test: NewStatusError(nil) -> [Invalid Content [http error]]
	//test: NewStatusError(nil) -> [Bad Request [http error]]

}

/*
type Request[T template.ErrorHandler] interface {
	Create(e T, req *http.Request) *http.Request
}

type Facebook struct{}

func (Facebook) Create(e template.ErrorHandler, req *http.Request) *http.Request {
	if e != nil {
	}

	return req
}

type Function interface {
	//Call(func() func(ctx context.Context, req *http.Request) *http.Request) *http.Request
	Call(req *http.Request) *http.Request
}

type TwitterRequest struct {
	Data string
}

func (d TwitterRequest) Call(req *http.Request) *http.Request {
	r, _ := http.NewRequest("GET", "http.www.google.com", nil)
	return r
}

func ExampleStatus_TemplateParameter() {

	//fmt.Printf("test: testAddress() -> %v\n", len(testAddress[address](address{})))
	//t := &TwitterRequest{}
	status, _ := do[template.DebugError, Facebook, []byte](nil)
	fmt.Printf("test: do() -> %v\n", status)
}

func do[E template.ErrorHandler, R Request[template.ErrorHandler], T any](req *http.Request) (*Status, T) {
	var t T
	var e E
	//var f F
	//f.Call(req)
	var r R

	r.Create(e, nil)
	return nil, t
}

func testAddress[T *address](param T) T {
	var t T
	var t2 address

	//t.GetZip()
	//param.State
	t2.GetZip()
	switch a := any(t).(type) {
	case address:
		return &a
	}
	return t //param
}

func testInt[T *int](param T) T {
	var t T

	*t += 6
	*param += 7

	switch a := any(t).(type) {
	case int:
		return &a
	}
	return param
}

func testString[T string](param T) T {
	var t T

	t += "6"
	param += "7"

	switch a := any(t).(type) {
	case string:
		return T(a)
	}
	return t
}

func testMap[T map[string]string](param T) T {
	var t T

	param["test"] = "first"
	t["test"] = "next"

	switch a := any(t).(type) {
	case map[string]string:
		a["test"] = "data"
		return a
	}
	return t
}

func testFunc[T func() int](param T) T {
	var t T

	param = func() int { return 0 }
	t = func() int { return 1 }

	switch a := any(t).(type) {
	case func() int:

		return a
	}
	return t
}


*/
