package io2

import (
	"encoding/json"
	"fmt"
)

func ExampleStatus_Marshal() {
	status := serializedStatusState{
		Code:     504,
		Location: "ExampleStatus2_Marshalling",
		Err:      "error 1",
	}
	s := ""
	buf, err := json.Marshal(status)
	if len(buf) > 0 {
		s = string(buf)
	}
	fmt.Printf("test: Marshal() -> [err:%v] [str:%v]\n", err, s)

	//Output:
	//test: Marshal() -> [err:<nil>] [str:{"code":504,"location":"ExampleStatus2_Marshalling","err":"error 1"}]

}
func ExampleNewS_Error() {
	status := NewStatusFrom("")
	fmt.Printf("test: NewStatusFrom(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	status = NewStatusFrom(s)
	fmt.Printf("test: NewStatusFrom(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	status = NewStatusFrom(s)
	fmt.Printf("test: NewStatusFrom(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: NewStatusFrom("") -> [status:OK]
	//test: NewStatusFrom(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: NewStatusFrom(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNewStatusFrom() {
	uri := "file://[cwd]/io2test/status-504.json"

	status := NewStatusFrom(uri)
	fmt.Printf("test: NewStatusFrom() -> [code:%v] [location:%v] [errors:[%v]]\n", status.Code, status.Trace(), status.Error())

	//Output:
	//test: NewStatusFrom() -> [code:504] [location:[ExampleStatus2_Marshalling]] [errors:[error 1]]

}

func ExampleNewStatusFrom_Const() {
	status := NewStatusFrom("")
	fmt.Printf("test: NewStatusFrom(nil) -> [code:%v]\n", status.Code)

	uri := StatusOKUri
	status = NewStatusFrom(uri)
	fmt.Printf("test: NewStatusFrom(\"%v\") -> [code:%v]\n", uri, status.Code)

	uri = StatusNotFoundUri
	status = NewStatusFrom(uri)
	fmt.Printf("test: NewStatusFrom(\"%v\") -> [code:%v] [status:%v]\n", uri, status.Code, status)

	uri = StatusTimeoutUri
	status = NewStatusFrom(uri)
	fmt.Printf("test: NewStatusFrom(\"%v\") -> [code:%v] [status:%v]\n", uri, status.Code, status)

	//Output:
	//test: NewStatusFrom(nil) -> [code:200]
	//test: NewStatusFrom("urn:status:ok") -> [code:200]
	//test: NewStatusFrom("urn:status:notfound") -> [code:404] [status:Not Found]
	//test: NewStatusFrom("urn:status:timeout") -> [code:504] [status:Timeout]

}
