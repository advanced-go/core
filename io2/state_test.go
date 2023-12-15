package io2

import (
	"fmt"
)

const (
	addressUrl   = "file://[cwd]/io2test/resource/address.json"
	status504Url = "file://[cwd]/io2test/resource/status-504.json"
)

type address struct {
	City    string
	State   string
	ZipCode string
}

func ExampleReadState() {
	t, status := ReadState[address](addressUrl)
	fmt.Printf("test: ReadState() -> [address:%v] [status:%v]\n", t, status)

	//Output:
	//test: ReadState() -> [address:{frisco texas 75034}] [status:OK]

}

func ExampleReadResults() {
	t, status := ReadResults[address](nil)
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	t, status = ReadResults[address]([]string{""})
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	t, status = ReadResults[address]([]string{"", ""})
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	t, status = ReadResults[address]([]string{addressUrl})
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	t, status = ReadResults[address]([]string{addressUrl, ""})
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	t, status = ReadResults[address]([]string{addressUrl, status504Url})
	fmt.Printf("test: ReadResults(nil) -> [state:%v] [status:%v]\n", t, status)

	//Output:
	//test: ReadResults(nil) -> [state:{  }] [status:OK]
	//test: ReadResults(nil) -> [state:{  }] [status:OK]
	//test: ReadResults(nil) -> [state:{  }] [status:OK]
	//test: ReadResults(nil) -> [state:{frisco texas 75034}] [status:OK]
	//test: ReadResults(nil) -> [state:{frisco texas 75034}] [status:OK]
	//test: ReadResults(nil) -> [state:{frisco texas 75034}] [status:Timeout [error 1]]

}
