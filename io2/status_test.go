package io2

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func ExampleReadStatus_Marshal() {
	status := statusState2{
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

func ExampleReadStatus_Unarshal() {
	uri, _ := url.Parse("file://[cwd]/io2test/resource/status-504.json")

	status := ReadStatus(uri)
	fmt.Printf("test: Unmarshal() -> [code:%v] [location:%v] [errors:%v]\n", status.Code(), status.Location(), status.Errors())

	//Output:
	//test: Unmarshal() -> [code:504] [location:[ExampleStatus2_Marshalling]] [errors:[error 1]]

}

func ExampleReadStatus_OK() {
	uri, _ := url.Parse(StatusOK)

	status := ReadStatus(uri)
	fmt.Printf("test: Unmarshal() -> [code:%v]\n", status.Code())

	//Output:
	//test: Unmarshal() -> [code:200]

}
