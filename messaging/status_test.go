package messaging

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"time"
)

// Example of how to get the value of an anonymous field in a struct
//
// https://go.dev/play/p/yQULMVaQK0

func ExampleNewStatus() {
	s := NewStatus(http.StatusTeapot)

	fmt.Printf("test: NewStatus() -> [status:%v] [ok:%v]\n", s, s.OK())

	//Output:
	//test: NewStatus() -> [status:I'm A Teapot] [ok:false]

}

func ExampleNewStatusError() {
	s := NewStatusError(http.StatusBadRequest, errors.New("bad request error"))

	fmt.Printf("test: NewStatusError() -> [status:%v] [ok:%v]\n", s, s.OK())
	var e runtime.Output

	e.Handle(s.Runtime(), "123-45-678")

	//Output:
	//test: NewStatusError() -> [status:Bad Request [bad request error]] [ok:false]
	//{ "code":400, "status":"Bad Request", "request-id":"123-45-678", "errors" : [ "bad request error" ], "trace" : [ "https://github.com/advanced-go/core/tree/main/messaging#ExampleNewStatusError","https://github.com/advanced-go/core/tree/main/messaging#NewStatusError" ] }

}

func ExampleNewStatusDuration() {
	s := NewStatusDuration(http.StatusOK, time.Millisecond*200)

	fmt.Printf("test: NewStatusDuration() -> [status:%v] [ok:%v] [duration:%v]\n", s, s.OK(), s.Duration)
	var e runtime.Output

	e.Handle(s.Runtime(), "123-45-678")

	//Output:
	//test: NewStatusDuration() -> [status:OK] [ok:true] [duration:200ms]

}
