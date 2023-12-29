package runtime

import (
	"fmt"
	"github.com/advanced-go/core/uri"
	"io"
	"os"
	"strings"
)

const (
	address1Url = "file://[cwd]/runtimetest/address1.json"
	address2Url = "file://[cwd]/runtimetest/address2.json"
)

type address2 struct {
	City    string
	State   string
	ZipCode string
}

func ExampleNew_StringError() {
	_, status := New[address2]("")
	fmt.Printf("test: New(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	_, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	_, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: New("") -> [status:Invalid Argument [error: URI is empty]]
	//test: New(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: New(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNew_URLError() {
	_, status := New[address2](nil)
	fmt.Printf("test: New(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	_, status = New[address2](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	_, status = New[address2](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: New("") -> [status:Invalid Argument [error: invalid type [<nil>]]]
	//test: New(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: New(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNew() {
	addr, status := New[address2](address1Url)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", address1Url, addr, status)

	addr, status = New[address2](uri.ParseRaw(address1Url))
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", address1Url, addr, status)

	buf, err := os.ReadFile(uri.FileName(address2Url))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	body := io.NopCloser(strings.NewReader(string(buf)))
	addr, status = New[address2](body)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", address2Url, addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/address2.json) -> [addr:{vinton iowa 52349}] [status:OK]

}
