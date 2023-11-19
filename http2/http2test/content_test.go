package http2test

import (
	"fmt"
	"github.com/advanced-go/core/http2"
	"github.com/advanced-go/core/io2"
	"net/http"
)

func Example_ReadContent_Empty() {
	s := "file://[cwd]/resource/get-request.txt"
	buf, err := io2.ReadFile(ParseRaw(s))
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
	buf, err := io2.ReadFile(ParseRaw(s))
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

func Example_ReadContentFromLocation() {
	h := make(http.Header)
	h.Add(http2.ContentLocation, "file://[cwd]/resource/activity.json")
	buf, status := ReadContentFromLocation(h)

	fmt.Printf("test: ReadContentFromLocation() -> [status:%v] [content:%v]\n", status, len(buf))

	//Output:
	//test: ReadContentFromLocation() -> [status:OK] [content:525]
	
}
