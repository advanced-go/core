package io2test

import (
	"encoding/json"
	"fmt"
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

func ExampleReadStatus_Unmarshal() {
	uri := "file://[cwd]/resource/status-504.json"

	status := ReadStatus(uri)
	fmt.Printf("test: Unmarshal() -> [code:%v] [location:%v] [errors:%v]\n", status.Code(), status.Location(), status.Errors())

	//Output:
	//test: Unmarshal() -> [code:504] [location:[ExampleStatus2_Marshalling]] [errors:[error 1]]

}

func ExampleReadStatus_OK() {
	status := ReadStatus(nil)
	fmt.Printf("test: ReadStatus(nil) -> [code:%v]\n", status.Code())

	uri := StatusOKUri
	status = ReadStatus(uri)
	fmt.Printf("test: ReadStatus(\"%v\") -> [code:%v]\n", uri, status.Code())

	uri2 := []string{""}
	status = ReadStatus(uri2)
	fmt.Printf("test: ReadStatus(\"%v\") -> [code:%v] [status:%v]\n", uri2, status.Code(), status)

	uri2 = []string{"", ""}
	status = ReadStatus(uri2)
	fmt.Printf("test: ReadStatus(\"%v\") -> [code:%v] [status:%v]\n", uri2, status.Code(), status)

	//Output:
	//test: ReadStatus(nil) -> [code:200]
	//test: ReadStatus("urn:status:ok") -> [code:200]
	//test: ReadStatus("[]") -> [code:200] [status:OK]
	//test: ReadStatus("[ ]") -> [code:200] [status:OK]

}
