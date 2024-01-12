package uri

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

func Example_FileNameError() {
	//s := "file://[cwd]/uritest/html-response.txt"
	//u, err := url.Parse(s)
	//fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	var t any
	name := FileName2(t)
	fmt.Printf("test: FileName2(nil) -> [type:%v] [url:%v]\n", reflect.TypeOf(t), name)

	s := ""
	name = FileName2(s)
	fmt.Printf("test: FileName2(\"\") -> [type:%v] [url:%v]\n", reflect.TypeOf(s), name)

	s = "https://www.google.com/search?q=golang"
	name = FileName2(s)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	s = "https://www.google.com/search?q=golang"
	u := parseRaw(s)
	name = FileName2(u)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	req, _ := http.NewRequest("", s, nil)
	name = FileName2(req)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(req), name)

	s = "file://[cwd]/uritest/html-response.txt"
	req, _ = http.NewRequest("", s, nil)
	name = FileName2(req)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(req), name)

	//Output:
	//test: FileName2(nil) -> [type:<nil>] [url:error: URL is nil]
	//test: FileName2("") -> [type:string] [url:error: URL is empty]
	//test: FileName2(https://www.google.com/search?q=golang) -> [type:string] [url:error: scheme is invalid [https]]
	//test: FileName2(https://www.google.com/search?q=golang) -> [type:*url.URL] [url:error: scheme is invalid [https]]
	//test: FileName2(https://www.google.com/search?q=golang) -> [type:*http.Request] [url:error: invalid URL type: *http.Request]
	//test: FileName2(file://[cwd]/uritest/html-response.txt) -> [type:*http.Request] [url:error: invalid URL type: *http.Request]

}

func Example_FileName() {
	s := "file://[cwd]/uritest/html-response.txt"
	u, err := url.Parse(s)
	fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	name := FileName2(s)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	name = FileName2(u)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	s = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	name = FileName2(s)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(s), name)

	u, err = url.Parse(s)
	name = FileName2(u)
	fmt.Printf("test: FileName2(%v) -> [type:%v] [url:%v]\n", s, reflect.TypeOf(u), name)

	//Output:
	//test: url.Parse(file://[cwd]/uritest/html-response.txt) -> [err:<nil>]
	//test: FileName2(file://[cwd]/uritest/html-response.txt) -> [type:string] [url:C:\Users\markb\GitHub\core\uri\uritest\html-response.txt]
	//test: FileName2(file://[cwd]/uritest/html-response.txt) -> [type:*url.URL] [url:C:\Users\markb\GitHub\core\uri\uritest\html-response.txt]
	//test: FileName2(file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt) -> [type:string] [url:c:\Users\markb\GitHub\core\uri\uritest\html-response.txt]
	//test: FileName2(file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt) -> [type:*url.URL] [url:c:\Users\markb\GitHub\core\uri\uritest\html-response.txt]

}

func Example_ReadFile() {
	s := "file://[cwd]/uritest/html-response.txt"
	u, _ := url.Parse(s)
	buf, err := os.ReadFile(FileName2(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	u, _ = url.Parse(s)
	buf, err = os.ReadFile(FileName2(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/uritest/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt) -> [err:<nil>] [buf:188]

}

/*
func Example_IsStatusURL() {
	u := ""
	fmt.Printf("test: IsStatusURL(\"%v\") -> %v\n", u, IsStatusURL(u))

	u = "file://[cwd]/io2test/resource/activity.json"
	fmt.Printf("test: IsStatusURL(\"%v\") -> %v\n", u, IsStatusURL(u))

	u = "file://[cwd]/io2test/resource/status/activity.json"
	fmt.Printf("test: IsStatusURL(\"%v\") -> %v\n", u, IsStatusURL(u))

	u = "file://[cwd]/io2test/resource/status-504.json"
	fmt.Printf("test: IsStatusURL(\"%v\") -> %v\n", u, IsStatusURL(u))

	u = "file://[cwd]/io2test/resource/status/status-504.json"
	fmt.Printf("test: IsStatusURL(\"%v\") -> %v\n", u, IsStatusURL(u))

	//Output:
	//test: IsStatusURL("") -> false
	//test: IsStatusURL("file://[cwd]/io2test/resource/activity.json") -> false
	//test: IsStatusURL("file://[cwd]/io2test/resource/status/activity.json") -> false
	//test: IsStatusURL("file://[cwd]/io2test/resource/status-504.json") -> true
	//test: IsStatusURL("file://[cwd]/io2test/resource/status/status-504.json") -> true

}


*/
