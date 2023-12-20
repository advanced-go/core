package uri

import (
	"fmt"
	"net/url"
	"os"
)

func Example_FileName() {
	s := "file://[cwd]/uritest/html-response.txt"
	u, err := url.Parse(s)

	fmt.Printf("test: url.Parse(%v) -> [err:%v]\n", s, err)

	name := FileName(u)
	fmt.Printf("test: FileName(%v) -> %v\n", s, name)

	s = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	u, err = url.Parse(s)
	name = FileName(u)
	fmt.Printf("test: FileName(%v) -> %v\n", s, name)

	//Output:
	//test: url.Parse(file://[cwd]/uritest/html-response.txt) -> [err:<nil>]
	//test: FileName(file://[cwd]/uritest/html-response.txt) -> C:\Users\markb\GitHub\core\uri\uritest\html-response.txt
	//test: FileName(file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt) -> c:\Users\markb\GitHub\core\uri\uritest\html-response.txt

}

func Example_ReadFile() {
	s := "file://[cwd]/uritest/html-response.txt"
	u, _ := url.Parse(s)
	buf, err := os.ReadFile(FileName(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	s = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	u, _ = url.Parse(s)
	buf, err = os.ReadFile(FileName(u))
	fmt.Printf("test: ReadFile(%v) -> [err:%v] [buf:%v]\n", s, err, len(buf))

	//Output:
	//test: ReadFile(file://[cwd]/uritest/html-response.txt) -> [err:<nil>] [buf:188]
	//test: ReadFile(file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt) -> [err:<nil>] [buf:188]

}