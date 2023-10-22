package exchange

import (
	"encoding/json"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"time"
)

type entryTest struct {
	Traffic    string
	Duration   time.Duration
	Region     string
	Zone       string
	SubZone    string
	Service    string
	Url        string
	Protocol   string
	Host       string
	Path       string
	Method     string
	StatusCode int32
}

func _Example_ReadRequest_GET() {
	s := "file://[cwd]/httptest/resource/http/get-request.txt"
	req, err := ReadRequest(runtime.ParseRaw(s))

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)

	//Output:
	//test: ReadRequest(file://[cwd]/httptest/resource/http/get-request.txt) -> [err:<nil>] [&{GET / HTTP/1.1 1 1 map[] {} <nil> 0 [] false foo.com map[] map[] <nil> map[]  / <nil> <nil> <nil> <nil>}]

}

func _Example_ReadRequest_Baseline() {
	s := "file://[cwd]/httptest/resource/http/baseline-request.txt"
	req, err := ReadRequest(runtime.ParseRaw(s))

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)

	//Output:
	//test: ReadRequest(file://[cwd]/httptest/resource/http/baseline-request.txt) -> [err:<nil>] [&{GET http://www.techcrunch.com/ HTTP/1.1 1 1 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8] Accept-Charset:[ISO-8859-1,utf-8;q=0.7,*;q=0.7] Accept-Encoding:[gzip,deflate\r\n] Accept-Language:[en-us,en;q=0.5] Content-Length:[7] Keep-Alive:[300] Proxy-Connection:[keep-alive] User-Agent:[Fake]] 0xc0000c6000 <nil> 7 [] false www.techcrunch.com map[] map[] <nil> map[]  http://www.techcrunch.com/ <nil> <nil> <nil> <nil>}]

}

func Example_ReadRequest_PUT() {
	s := "file://[cwd]/httptest/resource/http/test-put-req.txt"
	req, _ := ReadRequest(runtime.ParseRaw(s))

	buf, err1 := ReadAll[runtime.DebugError](req.Body)
	if err1 != nil {
	}
	var entry []entryTest
	json.Unmarshal(buf, &entry)
	fmt.Printf("test: ReadRequest(%v) -> [cnt:%v] [fields:%v]\n", s, len(entry), entry)

	//Output:
	//test: ReadRequest(file://[cwd]/httptest/resource/http/test-put-req.txt) -> [cnt:1] [fields:[{ingress 800µs usa west  access-log https://access-log.com/example-domain/timeseries/entry http access-log.com /example-domain/timeseries/entry GET 200}]]

}
