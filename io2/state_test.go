package io2

import (
	"fmt"
	"net/url"
)

type address struct {
	City    string
	State   string
	ZipCode string
}

func ExampleReadState() {
	u, _ := url.Parse("file://[cwd]/io2test/resource/address.json")
	t, status := ReadState[address](u)
	fmt.Printf("test: ReadState() -> [address:%v] [status:%v]\n", t, status)

	//Output:
	//test: ReadState() -> [address:{frisco texas 75034}] [status:OK]
	
}
