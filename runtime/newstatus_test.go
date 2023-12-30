package runtime

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
	status := NewS("")
	fmt.Printf("test: NewS(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	status = NewS(s)
	fmt.Printf("test: NewS(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	status = NewS(s)
	fmt.Printf("test: NewS(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: NewS("") -> [status:OK]
	//test: NewS(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: NewS(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNewS() {
	uri := "file://[cwd]/runtimetest/status-504.json"

	status := NewS(uri)
	fmt.Printf("test: NewS() -> [code:%v] [location:%v] [errors:%v]\n", status.Code(), status.Location(), status.Errors())

	//Output:
	//test: NewS() -> [code:504] [location:[ExampleStatus2_Marshalling]] [errors:[error 1]]

}

func ExampleNewS_Const() {
	status := NewS("")
	fmt.Printf("test: NewS(nil) -> [code:%v]\n", status.Code())

	uri := StatusOKUri
	status = NewS(uri)
	fmt.Printf("test: NewS(\"%v\") -> [code:%v]\n", uri, status.Code())

	uri = StatusNotFoundUri
	status = NewS(uri)
	fmt.Printf("test: NewS(\"%v\") -> [code:%v] [status:%v]\n", uri, status.Code(), status)

	uri = StatusTimeoutUri
	status = NewS(uri)
	fmt.Printf("test: NewS(\"%v\") -> [code:%v] [status:%v]\n", uri, status.Code(), status)

	//Output:
	//test: NewS(nil) -> [code:200]
	//test: NewS("urn:status:ok") -> [code:200]
	//test: NewS("urn:status:notfound") -> [code:404] [status:Not Found]
	//test: NewS("urn:status:timeout") -> [code:504] [status:Timeout]

}
