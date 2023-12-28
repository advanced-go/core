package http2test

import (
	"encoding/json"
	"fmt"
	"github.com/advanced-go/core/http2/io2"
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

func Example_ReadRequest_GET() {
	s := "file://[cwd]/resource/get-request.txt"
	req, err := ReadRequest(ParseRaw(s))

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)

	//Output:
	//test: ReadRequest(file://[cwd]/resource/get-request.txt) -> [err:<nil>] [&{GET /advanced-go/example-domain/activity/entry?v=v2 HTTP/1.1 1 1 map[Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8] Accept-Charset:[ISO-8859-1,utf-8;q=0.7,*;q=0.7] Accept-Encoding:[gzip,deflate\r\n] Accept-Language:[en-us,en;q=0.5] Connection:[close] Content-Length:[7] Content-Location:[github.com/advanced-go/example-domain/activity/EntryV1] Keep-Alive:[300] Proxy-Connection:[keep-alive] User-Agent:[Fake]] {
	//} <nil> 7 [] true localhost:8080 map[] map[] <nil> map[]  /advanced-go/example-domain/activity/entry?v=v2 <nil> <nil> <nil> <nil>}]

}

func Example_ReadRequest_Baseline() {
	s := "file://[cwd]/resource/baseline-request.txt"
	req, err := ReadRequest(ParseRaw(s))

	if req != nil {
	}
	// print content
	//fmt.Printf("test: ReadRequest(%v) -> [err:%v] [%v]\n", s, err, req)
	fmt.Printf("test: ReadRequest(%v) -> [err:%v]\n", s, err)

	//Output:
	//test: ReadRequest(file://[cwd]/resource/baseline-request.txt) -> [err:<nil>]

}

func Example_ReadRequest_PUT() {
	s := "file://[cwd]/resource/put-req.txt"
	req, err := ReadRequest(ParseRaw(s))

	if err != nil {
		fmt.Printf("test: ReadRequest(%v) -> [err:%v]\n", s, err)
	} else {
		buf, err1 := io2.ReadAll(req.Body)
		if err1 != nil {
		}
		var entry []entryTest
		json.Unmarshal(buf, &entry)
		fmt.Printf("test: ReadRequest(%v) -> [cnt:%v] [fields:%v]\n", s, len(entry), entry)
	}

	//Output:
	//test: ReadRequest(file://[cwd]/resource/put-req.txt) -> [cnt:2] [fields:[{ingress 800µs usa west  access-log https://access-log.com/example-domain/timeseries/entry http access-log.com /example-domain/timeseries/entry GET 200} {egress 100µs usa east  access-log https://access-log.com/example-domain/timeseries/entry http access-log.com /example-domain/timeseries/entry PUT 202}]]

}
