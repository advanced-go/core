package io2

import (
	"fmt"
)

type address struct {
	City    string
	State   string
	ZipCode string
}

func ExampleReadState() {
	uri := "file://[cwd]/io2test/resource/address.json"
	t, status := ReadState[address](uri)
	fmt.Printf("test: ReadState() -> [address:%v] [status:%v]\n", t, status)

	//Output:
	//test: ReadState() -> [address:{frisco texas 75034}] [status:OK]

}
