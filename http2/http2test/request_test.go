package http2test

import (
	"encoding/json"
	"fmt"
	"github.com/advanced-go/core/io2"
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

	//buf,status :=

	fmt.Printf("test: ReadRequest(%v) -> [err:%v] [content-location:%v]\n", s, err, req.Header.Get("Content-Location"))

	//Output:
	//test: ReadRequest(file://[cwd]/resource/get-request.txt) -> [err:<nil>] [content-location:github/advanced-go/example-domain/activity/EntryV1]

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
		buf, err1 := io2.ReadAll(req.Body, nil)
		if err1 != nil {
		}
		var entry []entryTest
		json.Unmarshal(buf, &entry)
		fmt.Printf("test: ReadRequest(%v) -> [cnt:%v] [fields:%v]\n", s, len(entry), entry)
	}

	//Output:
	//test: ReadRequest(file://[cwd]/resource/put-req.txt) -> [cnt:2] [fields:[{ingress 800µs usa west  access-log https://access-log.com/example-domain/timeseries/entry http access-log.com /example-domain/timeseries/entry GET 200} {egress 100µs usa east  access-log https://access-log.com/example-domain/timeseries/entry http access-log.com /example-domain/timeseries/entry PUT 202}]]

}
