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

func ExampleNewBytes_Uri() {
	s := status504
	buf, status := newBytes(s)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = address1Url
	buf, status = newBytes(s)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = status504
	u := uri.ParseRaw(s)
	buf, status = newBytes(u)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	s = address1Url
	u = uri.ParseRaw(s)
	buf, status = newBytes(u)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	//Output:
	//test: newBytes(file://[cwd]/runtimetest/status-504.json) -> [type:string] [buf:82] [status:OK]
	//test: newBytes(file://[cwd]/runtimetest/address1.json) -> [type:string] [buf:68] [status:OK]
	//test: newBytes(file://[cwd]/runtimetest/status-504.json) -> [type:*url.URL] [buf:82] [status:OK]
	//test: newBytes(file://[cwd]/runtimetest/address1.json) -> [type:*url.URL] [buf:68] [status:OK]

}

func ExampleNewBytes_Bytes() {
	s := address2Url
	buf, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	buf1, status := newBytes(buf)
	fmt.Printf("test: newBytes(%v) -> [in:%v] [out:%v] [status:%v]\n", s, len(buf), len(buf1), status)

	//Output:
	//test: newBytes(file://[cwd]/runtimetest/address2.json) -> [in:67] [out:67] [status:OK]

}

func ExampleNewBytes_Reader() {
	s := address3Url
	buf0, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}

	r := strings.NewReader(string(buf0))
	buf, status := newBytes(r)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(r), len(buf), status)

	body := io.NopCloser(strings.NewReader(string(buf0)))
	buf, status = newBytes(body)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(body), len(buf), status)

	//Output:
	//test: newBytes(file://[cwd]/runtimetest/address3.json) -> [type:*strings.Reader] [buf:72] [status:OK]
	//test: newBytes(file://[cwd]/runtimetest/address3.json) -> [type:io.nopCloserWriterTo] [buf:72] [status:OK]

}

func ExampleNewBytes_Response() {
	s := address3Url
	buf0, err := os.ReadFile(uri.FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := new(http.Response)
	r.Body = io.NopCloser(strings.NewReader(string(buf0)))

	buf, status := newBytes(r)
	fmt.Printf("test: newBytes(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(r), len(buf), status)

	//Output:
	//test: newBytes(file://[cwd]/runtimetest/address3.json) -> [type:*http.Response] [buf:72] [status:OK]

}
