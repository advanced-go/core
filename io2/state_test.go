package io2

import (
	"fmt"
)

const (
	addressUrl   = "file://[cwd]/io2test/resource/address1.json"
	status504Url = "file://[cwd]/io2test/resource/status-504.json"
)

type address struct {
	City    string
	State   string
	ZipCode string
}

func ExampleReadState_Error() {
	_, status := ReadState[address]("")
	fmt.Printf("test: ReadState(\"\") -> [status:%v]\n", status)

	//var list []string
	//_, status = ReadState[runtime.Nillable](list)
	//fmt.Printf("test: ReadState(%v) -> [status:%v]\n", list, status)

	//list = []string{"", ""}
	//_, status = ReadState[runtime.Nillable](list)
	//fmt.Printf("test: ReadState(%v) -> [status:%v]\n", list, status)

	//n := 1234
	//_, status = ReadState[runtime.Nillable](n)
	//fmt.Printf("test: ReadState(%v) -> [status:%v]\n", n, status)

	//Output:
	//test: ReadState("") -> [status:Invalid Argument [error: URI is empty]]

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
