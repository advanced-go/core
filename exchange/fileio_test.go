package exchange

import (
	"fmt"
	"net/url"
)

func Example_ReadFile() {
	s := "file://[cwd]/httptest/resource/http/html-response.html"
	buf, err := ReadFile(s)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html"
	buf, err = ReadFile(s)
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/exchange/httptest/resource/http/html-response.html) -> [err:<nil>] [buf:188]

}

func Example_createFName() {
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
