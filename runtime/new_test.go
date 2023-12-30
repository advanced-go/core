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
	status504   = "file://[cwd]/runtimetest/status-504.json"
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
	s := status504
	addr, status := New[address2](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = address1Url
	addr, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = status504
	addr, status = New[address2](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = address1Url
	addr, status = New[address2](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	body := io.NopCloser(strings.NewReader(string(buf)))
	addr, status = New[address2](body)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/address2.json) -> [addr:{vinton iowa 52349}] [status:OK]

}

func ExampleNew_Const_Status() {
	s := StatusOKUri
	addr, status := New[address2](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = StatusNotFoundUri
	addr, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = StatusTimeoutUri
	addr, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	//Output:
	//test: New(urn:status:ok) -> [addr:{  }] [status:OK]
	//test: New(urn:status:notfound) -> [addr:{  }] [status:Not Found]
	//test: New(urn:status:timeout) -> [addr:{  }] [status:Timeout]

}
