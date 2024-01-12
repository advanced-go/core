package runtime

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

func ExampleNewBytes_Uri() {
	s := status504
	buf, status := NewBytes(s)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = address1Url
	buf, status = NewBytes(s)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = status504
	u := parseRaw(s)
	buf, status = NewBytes(u)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	s = address1Url
	u = parseRaw(s)
	buf, status = NewBytes(u)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	//Output:
	//test: NewBytes(file://[cwd]/runtimetest/status-504.json) -> [type:string] [buf:82] [status:OK]
	//test: NewBytes(file://[cwd]/runtimetest/address1.json) -> [type:string] [buf:68] [status:OK]
	//test: NewBytes(file://[cwd]/runtimetest/status-504.json) -> [type:*url.URL] [buf:82] [status:OK]
	//test: NewBytes(file://[cwd]/runtimetest/address1.json) -> [type:*url.URL] [buf:68] [status:OK]

}

func ExampleNewBytes_Bytes() {
	s := address2Url
	buf, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	buf1, status := NewBytes(buf)
	fmt.Printf("test: NewBytes(%v) -> [in:%v] [out:%v] [status:%v]\n", s, len(buf), len(buf1), status)

	//Output:
	//test: NewBytes(file://[cwd]/runtimetest/address2.json) -> [in:67] [out:67] [status:OK]

}

func ExampleNewBytes_Reader() {
	s := address3Url
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}

	r := strings.NewReader(string(buf0))
	buf, status := NewBytes(r)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(r), len(buf), status)

	body := io.NopCloser(strings.NewReader(string(buf0)))
	buf, status = NewBytes(body)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(body), len(buf), status)

	//Output:
	//test: NewBytes(file://[cwd]/runtimetest/address3.json) -> [type:*strings.Reader] [buf:72] [status:OK]
	//test: NewBytes(file://[cwd]/runtimetest/address3.json) -> [type:io.nopCloserWriterTo] [buf:72] [status:OK]

}

func ExampleNewBytes_Response() {
	s := address3Url
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := new(http.Response)
	r.Body = io.NopCloser(strings.NewReader(string(buf0)))

	buf, status := NewBytes(r)
	fmt.Printf("test: NewBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(r), len(buf), status)

	//Output:
	//test: NewBytes(file://[cwd]/runtimetest/address3.json) -> [type:*http.Response] [buf:72] [status:OK]

}
