package io

import (
	"fmt"
	"net/url"
)

// ParseRaw - parse a raw Uri without error
func ParseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

func Example_ReadFileRaw() {
	s := "file://[cwd]/iotest/resource/html-response.txt"
	buf, err := ReadFile(ParseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/io/iotest/resource/html-response.txt"
	buf, err = ReadFile(ParseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFileRaw(file://[cwd]/iotest/resource/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFileRaw(file:///c:/Users/markb/GitHub/core/io/iotest/resource/html-response.txt) -> [err:<nil>] [buf:188]

}

func Example_ReadFile() {
	s := "file://[cwd]/iotest/resource/html-response.txt"
	u, _ := url.Parse(s)
	buf, err := ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/io/iotest/resource/html-response.txt"
	u, _ = url.Parse(s)
	buf, err = ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/iotest/resource/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/io/iotest/resource/html-response.txt) -> [err:<nil>] [buf:188]

}

/*
func _Example_createFName() {
	s := "file://[cwd]/iotest/resource/http/html-response.txt"
	u, err := url.Parse(s)

	fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	name := createFname(u)
	fmt.Printf("test: createFname(%v) -> %v\n", s, name)

	s = "file:///c:/Users/markb/GitHub/core/io/iotest/resource/http/html-response.txt"
	u, err = url.Parse(s)
	name = createFname(u)
	fmt.Printf("test: createFname(%v) -> %v\n", s, name)

	//Output:
	//test: url.Parse(file://[cwd]/iotest/resource/http/html-response.txt) -> [err:<nil>]
	//test: createFname(file://[cwd]/iotest/resource/http/html-response.txt) -> C:\Users\markb\GitHub\core\io\iotest\resource\http\html-response.txt
	//test: createFname(file:///c:/Users/markb/GitHub/core/io/iotest/resource/http/html-response.txt) -> c:\Users\markb\GitHub\core\io\iotest\resource\http\html-response.txt

}


*/
