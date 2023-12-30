package runtime

import (
	"fmt"
	"github.com/advanced-go/core/uri"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	address1Url = "file://[cwd]/runtimetest/address1.json"
	address2Url = "file://[cwd]/runtimetest/address2.json"
	address3Url = "file://[cwd]/runtimetest/address3.json"
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

func ExampleNew_Uri() {
	s := status504
	addr, status := New[address2](s)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(s), addr, status)

	s = address1Url
	addr, status = New[address2](s)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(s), addr, status)

	s = status504
	u := uri.ParseRaw(s)
	addr, status = New[address2](u)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(u), addr, status)

	s = address1Url
	u = uri.ParseRaw(s)
	addr, status = New[address2](u)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(u), addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [type:string] [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [type:string] [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [type:*url.URL] [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [type:*url.URL] [addr:{frisco texas 75034}] [status:OK]

}

func ExampleNew_Bytes() {
	s := address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	addr, status := New[address2](buf)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(buf), addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/address2.json) -> [type:[]uint8] [addr:{vinton iowa 52349}] [status:OK]

}

func ExampleNew_Reader() {
	s := address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := strings.NewReader(string(buf))
	addr, status := New[address2](r)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(r), addr, status)

	s = address3Url
	buf, err = os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	body := io.NopCloser(strings.NewReader(string(buf)))
	addr, status = New[address2](body)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(body), addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/address2.json) -> [type:*strings.Reader] [addr:{vinton iowa 52349}] [status:OK]
	//test: New(file://[cwd]/runtimetest/address3.json) -> [type:io.nopCloserWriterTo] [addr:{forest city iowa 50436}] [status:OK]

}

func ExampleNew_Response() {
	s := address3Url
	buf0, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := new(http.Response)
	r.Body = io.NopCloser(strings.NewReader(string(buf0)))

	addr, status := New[address2](r)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [addr%v] [status:%v]\n", s, reflect.TypeOf(r), addr, status)

	//Output:
	//test: NewBytes(file://[cwd]/runtimetest/address3.json) -> [type:*http.Response] [addr{forest city iowa 50436}] [status:OK]

}
