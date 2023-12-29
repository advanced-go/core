package runtime

import (
	"encoding/json"
	"fmt"
)

const (
	addressUrl = "file://[cwd]/runtimetest/address.json"
	//status504Url = "file://[cwd]/io2test/resource/status-504.json"
)

type address2 struct {
	City    string
	State   string
	ZipCode string
}

func ExampleNewT_Error() {
	_, status := NewT[address2]("")
	fmt.Printf("test: NewT(\"\") -> [status:%v]\n", status)

	//Output:
	//test: NewT("") -> [status:Invalid Argument [error: URI is empty]]

}

func ExampleNewT() {
	t, status := NewT[address2](addressUrl)
	fmt.Printf("test: NewT() -> [address:%v] [status:%v]\n", t, status)

	//Output:
	//test: NewT() -> [address:{frisco texas 75034}] [status:OK]

}

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
