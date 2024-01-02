package runtime

import (
	"fmt"
	"github.com/advanced-go/core/uri"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	address1Url = "file://[cwd]/runtimetest/address1.json"
	address2Url = "file://[cwd]/runtimetest/address2.json"
	address3Url = "file://[cwd]/runtimetest/address3.json"
	status504   = "file://[cwd]/runtimetest/status-504.json"
)

type newAddress struct {
	City    string
	State   string
	ZipCode string
}

func ExampleNew_String_Error() {
	_, status := New[newAddress]("")
	fmt.Printf("test: New(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	_, status = New[newAddress](s)
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	_, status = New[newAddress](s)
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: New("") -> [status:Invalid Argument [error: URI is empty]]
	//test: New(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: New(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNew_String_Status() {
	s := StatusOKUri
	addr, status := New[newAddress](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	s = StatusNotFoundUri
	bytes, status0 := New[[]byte](s)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, bytes, status0)

	s = StatusTimeoutUri
	addr, status = New[newAddress](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	// status from uri, generic type is ignored
	s = status504
	bytes1, status1 := New[[]byte](s)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, bytes1, status1)

	//Output:
	//test: New(urn:status:ok) -> [addr:{  }] [status:OK]
	//test: New(urn:status:notfound) -> [bytes:[]] [status:Not Found]
	//test: New(urn:status:timeout) -> [addr:{  }] [status:Timeout]
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [bytes:[]] [status:Timeout [error 1]]

}

func ExampleNew_String_URI() {
	// bytes
	s := address1Url
	bytes, status := New[[]byte](s)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status)

	// type
	s = address1Url
	addr, status1 := New[address](s)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status1)

	//Output:
	//test: New(file://[cwd]/runtimetest/address1.json) -> [bytes:68] [status:OK]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{ frisco texas }] [status:OK]

}

func ExampleNew_URL_Error() {
	_, status := New[newAddress](nil)
	fmt.Printf("test: New(\"\") -> [status:%v]\n", status)

	s := "https://www.google.com/search"
	_, status = New[newAddress](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	s = "file://[cwd]/runtimetest/address.txt"
	_, status = New[newAddress](uri.ParseRaw(s))
	fmt.Printf("test: New(%v) -> [status:%v]\n", s, status)

	//Output:
	//test: New("") -> [status:Invalid Argument [error: invalid type [<nil>]]]
	//test: New(https://www.google.com/search) -> [status:Invalid Argument [error: URI is not of scheme file: https://www.google.com/search]]
	//test: New(file://[cwd]/runtimetest/address.txt) -> [status:Invalid Argument [error: URI is not a JSON file]]

}

func ExampleNew_URL_Status() {
	s := status504
	u, _ := url.Parse(s)

	addr, status0 := New[newAddress](u)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status0)

	//Output:
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [addr:{  }] [status:Timeout [error 1]]

}

func ExampleNew_URL() {
	s := address1Url
	u, _ := url.Parse(s)
	bytes, status0 := New[[]byte](u)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status0)

	s = address1Url
	u, _ = url.Parse(s)
	addr, status1 := New[newAddress](u)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status1)

	//Output:
	//test: New(file://[cwd]/runtimetest/address1.json) -> [bytes:68] [status:OK]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [addr:{frisco texas 75034}] [status:OK]

}

func ExampleNew_Bytes() {
	s := address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}

	bytes, status0 := New[[]byte](buf)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status0)

	addr, status := New[newAddress](buf)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/address2.json) -> [bytes:67] [status:OK]
	//test: New(file://[cwd]/runtimetest/address2.json) -> [addr:{vinton iowa 52349}] [status:<nil>]

}

func ExampleNew_Response() {
	s := address3Url
	buf0, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := new(http.Response)
	r.Body = io.NopCloser(strings.NewReader(string(buf0)))

	bytes, status := New[[]byte](r)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status)

	r = new(http.Response)
	r.Body = io.NopCloser(strings.NewReader(string(buf0)))
	addr, status1 := New[newAddress](r)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status1)

	//Output:
	//test: New(file://[cwd]/runtimetest/address3.json) -> [bytes:72] [status:OK]
	//test: New(file://[cwd]/runtimetest/address3.json) -> [addr:{forest city iowa 50436}] [status:OK]

}

/*

func ExampleNew_Uri() {
	s := status504
	addr, status := New[newAddress](s)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(s), addr, status)

	s = address1Url
	addr, status = New[newAddress](s)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(s), addr, status)

	s = status504
	u := uri.ParseRaw(s)
	addr, status = New[newAddress](u)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(u), addr, status)

	s = address1Url
	u = uri.ParseRaw(s)
	addr, status = New[newAddress](u)
	fmt.Printf("test: New(%v) -> [type:%v] [addr:%v] [status:%v]\n", s, reflect.TypeOf(u), addr, status)

	//Output:
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [type:string] [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [type:string] [addr:{frisco texas 75034}] [status:OK]
	//test: New(file://[cwd]/runtimetest/status-504.json) -> [type:*url.URL] [addr:{  }] [status:Timeout [error 1]]
	//test: New(file://[cwd]/runtimetest/address1.json) -> [type:*url.URL] [addr:{frisco texas 75034}] [status:OK]

}


*/

func ExampleNew_Reader() {
	s := address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := strings.NewReader(string(buf))
	bytes, status := New[[]byte](r)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status)

	r = strings.NewReader(string(buf))
	addr, status1 := New[newAddress](r)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status1)

	//Output:
	//test: New(file://[cwd]/runtimetest/address2.json) -> [bytes:67] [status:OK]
	//test: New(file://[cwd]/runtimetest/address2.json) -> [addr:{vinton iowa 52349}] [status:OK]

}

func ExampleNew_ReadCloser() {
	s := address3Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	body := io.NopCloser(strings.NewReader(string(buf)))
	bytes, status := New[[]byte](body)
	fmt.Printf("test: New(%v) -> [bytes:%v] [status:%v]\n", s, len(bytes), status)

	body = io.NopCloser(strings.NewReader(string(buf)))
	addr, status1 := New[newAddress](body)
	fmt.Printf("test: New(%v) -> [addr:%v] [status:%v]\n", s, addr, status1)

	//Output:
	//test: New(file://[cwd]/runtimetest/address3.json) -> [bytes:72] [status:OK]
	//test: New(file://[cwd]/runtimetest/address3.json) -> [addr:{forest city iowa 50436}] [status:OK]

}
