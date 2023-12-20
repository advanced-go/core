package io2

import (
	"fmt"
	"net/url"
)

// ParseRaw - parse a raw Uri without error
func parseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

func Example_ReadFileRaw() {
	s := "file://[cwd]/io2test/resource/html-response.txt"
	buf, err := ReadFile(parseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/io2/io2test/resource/html-response.txt"
	buf, err = ReadFile(parseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFileRaw(file://[cwd]/io2test/resource/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFileRaw(file:///c:/Users/markb/GitHub/core/io2/io2test/resource/html-response.txt) -> [err:<nil>] [buf:188]

}

func Example_ReadFile() {
	s := "file://[cwd]/io2test/resource/html-response.txt"
	u, _ := url.Parse(s)
	buf, err := ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/io2/io2test/resource/html-response.txt"
	u, _ = url.Parse(s)
	buf, err = ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/io2test/resource/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/io2/io2test/resource/html-response.txt) -> [err:<nil>] [buf:188]

}

func Example_ReadFileFromPath() {
	buf, status := ReadFileFromPath("file://[cwd]/io2test/resource/activity.json")
	fmt.Printf("test: ReadContentFromLocation() -> [status:%v] [content:%v]\n", status, len(buf))

	//Output:
	//test: ReadContentFromLocation() -> [status:OK] [content:525]

}
