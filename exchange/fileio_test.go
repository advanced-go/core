package exchange

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/url"
)

func Example_ReadFileRaw() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	buf, err := ReadFile(runtime.ParseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html"
	buf, err = ReadFile(runtime.ParseRaw(s))
	fmt.Printf("test: ReadFileRaw(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFileRaw(file://[cwd]/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]
	//test: ReadFileRaw(file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]

}

func Example_ReadFile() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	u, _ := url.Parse(s)
	buf, err := ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html"
	u, _ = url.Parse(s)
	buf, err = ReadFile(u)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]

}

/*
func _Example_createFName() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	u, err := url.Parse(s)

	fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	name := createFname(u)
	fmt.Printf("test: createFname(%v) -> %v\n", s, name)

	s = "file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html"
	u, err = url.Parse(s)
	name = createFname(u)
	fmt.Printf("test: createFname(%v) -> %v\n", s, name)

	//Output:
	//test: url.Parse(file://[cwd]/httptest/resource/http/html-response.html) -> [err:<nil>]
	//test: createFname(file://[cwd]/httptest/resource/http/html-response.html) -> C:\Users\markb\GitHub\core\exchange\httptest\resource\http\html-response.html
	//test: createFname(file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html) -> c:\Users\markb\GitHub\core\exchange\httptest\resource\http\html-response.html

}


*/
