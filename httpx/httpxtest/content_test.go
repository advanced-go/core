package httpxtest

import (
	"fmt"
	"github.com/go-ai-agent/core/httpx"
)

func Example_ReadContent_Empty() {
	s := "file://[cwd]/resource/get-request.txt"
	buf, err := httpx.ReadFile(httpx.ParseRaw(s))
	if err != nil {
		fmt.Printf("test: ReadFile(%v) -> [err:%v]\n", s, err)

	} else {
		bytes, err1 := ReadContent(buf)
		fmt.Printf("test: ReadContent() -> [err:%v] [bytes:%v]\n", err1, bytes.Len())
	}

	//Output:
	//test: ReadContent() -> [err:<nil>] [bytes:0]

}

func _Example_ReadContent_Available() {
	s := "file://[cwd]/resource/put-req.txt"
	buf, err := httpx.ReadFile(httpx.ParseRaw(s))
	if err != nil {
		fmt.Printf("test: ReadFile(%v) -> [err:%v]\n", s, err)

	} else {
		bytes, err1 := ReadContent(buf)
		fmt.Printf("test: ReadContent() -> [err:%v] [bytes:%v] %v\n", err1, bytes.Len(), bytes.String())
	}

	//Output:
	//test: ReadContent() -> [err:<nil>] [bytes:872] [
	//  {
	//    "Traffic":     "ingress",
	//    "Duration":    800000,
	//    "Region":      "usa",
	//    "Zone":        "west",
	//    "SubZone":     "",
	//    "Service":     "access-log",
	//    "Url":         "https://access-log.com/example-domain/timeseries/entry",
	//    "Protocol":    "http",
	//    "Host":        "access-log.com",
	//    "Path":        "/example-domain/timeseries/entry",
	//    "Method":      "GET",
	//    "StatusCode":  200
	//  },
	//  {
	//    "Traffic":     "egress",
	//    "Duration":    100000,
	//    "Region":      "usa",
	//    "Zone":        "east",
	//    "SubZone":     "",
	//    "Service":     "access-log",
	//    "Url":         "https://access-log.com/example-domain/timeseries/entry",
	//    "Protocol":    "http",
	//    "Host":        "access-log.com",
	//    "Path":        "/example-domain/timeseries/entry",
	//    "Method":      "PUT",
	//    "StatusCode":  202
	//  }
	//]

}
