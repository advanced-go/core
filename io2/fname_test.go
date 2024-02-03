package io2

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

// parseRaw - parse a raw Uri without error
func parseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}

func Example_FileNameError() {
	//s := "file://[cwd]/io2test/test-response.txt"
	//u, err := url.Parse(s)
	//fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	var t any
	name := FileName(t)
	fmt.Printf("test: FileName(nil) -> [type:%v] [url:%v]\n", reflect.TypeOf(t), name)

	s := ""
	name = FileName(s)
	fmt.Printf("test: FileName(\"\") -> [type:%v] [url:%v]\n", reflect.TypeOf(s), name)

	s = "https://www.google.com/search?q=golang"
	name = FileName(s)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	s = "https://www.google.com/search?q=golang"
	u := parseRaw(s)
	name = FileName(u)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	req, _ := http.NewRequest("", s, nil)
	name = FileName(req)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(req), name)

	s = "file://[cwd]/io2test/test-response.txt"
	req, _ = http.NewRequest("", s, nil)
	name = FileName(req)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(req), name)

	//Output:
	//test: FileName(nil) -> [type:<nil>] [url:error: URL is nil]
	//test: FileName("") -> [type:string] [url:error: URL is empty]
	//test: FileName(https://www.google.com/search?q=golang) -> [type:string] [url:error: scheme is invalid [https]]
	//test: FileName(https://www.google.com/search?q=golang) -> [type:*url.URL] [url:error: scheme is invalid [https]]
	//test: FileName(https://www.google.com/search?q=golang) -> [type:*http.Request] [url:error: invalid URL type: *http.Request]
	//test: FileName(file://[cwd]/io2test/test-response.txt) -> [type:*http.Request] [url:error: invalid URL type: *http.Request]

}

func Example_FileName() {
	s := "file://[cwd]/io2test/test-response.txt"
	u, err := url.Parse(s)
	fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	name := FileName(s)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	name = FileName(u)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	s = "file:///c:/Users/markb/GitHub/core/io2/io2test/test-response.txt"
	name = FileName(s)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	u, err = url.Parse(s)
	name = FileName(u)
	fmt.Printf("test: FileName(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	//Output:
	//test: url.Parse(file://[cwd]/io2test/test-response.txt) -> [err:<nil>]
	//test: FileName(file://[cwd]/io2test/test-response.txt) -> [type:string] [url:C:\Users\markb\GitHub\core\io2\io2test\test-response.txt]
	//test: FileName(file://[cwd]/io2test/test-response.txt) -> [type:*url.URL] [url:C:\Users\markb\GitHub\core\io2\io2test\test-response.txt]
	//test: FileName(file:///c:/Users/markb/GitHub/core/io2/io2test/test-response.txt) -> [type:string] [url:c:\Users\markb\GitHub\core\io2\io2test\test-response.txt]
	//test: FileName(file:///c:/Users/markb/GitHub/core/io2/io2test/test-response.txt) -> [type:*url.URL] [url:c:\Users\markb\GitHub\core\io2\io2test\test-response.txt]

}

func Example_ReadFile() {
	s := "file://[cwd]/io2test/test-response.txt"
	u, _ := url.Parse(s)
	buf, err := os.ReadFile(FileName(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/io2/io2test/test-response.txt"
	u, _ = url.Parse(s)
	buf, err = os.ReadFile(FileName(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/io2test/test-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/io2/io2test/test-response.txt) -> [err:<nil>] [buf:188]

}
