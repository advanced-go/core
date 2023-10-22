package exchange

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

func Example_ReadRequest_GET() {
	s := "file://[cwd]/httptest/resource/http/get-request.txt"
	req, err := ReadRequest(runtime.ParseRaw(s))

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)

	//Output:
	//test: ReadRequest(file://[cwd]/httptest/resource/http/get-request.txt) -> [err:<nil>] [&{GET / HTTP/1.1 1 1 map[] {} <nil> 0 [] false foo.com map[] map[] <nil> map[]  / <nil> <nil> <nil> <nil>}]

}

func Example_ReadRequest_Baseline() {
	s := "file://[cwd]/httptest/resource/http/baseline-request.txt"
	req, err := ReadRequest(runtime.ParseRaw(s))

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)

	//Output:
	//test: ReadRequest(file://[cwd]/httptest/resource/http/baseline-request.txt) -> [err:<nil>] [&{GET http://www.techcrunch.com/ HTTP/1.1 1 1 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8] Accept-Charset:[ISO-8859-1,utf-8;q=0.7,*;q=0.7] Accept-Encoding:[gzip,deflate\r\n] Accept-Language:[en-us,en;q=0.5] Content-Length:[7] Keep-Alive:[300] Proxy-Connection:[keep-alive] User-Agent:[Fake]] 0xc0000c6000 <nil> 7 [] false www.techcrunch.com map[] map[] <nil> map[]  http://www.techcrunch.com/ <nil> <nil> <nil> <nil>}]
	
}
