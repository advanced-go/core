package runtime

import (
	"fmt"
	"net/http"
	"time"
)

func _ExampleFormatLogJson() {
	start := time.Now().UTC()
	req, _ := http.NewRequest(http.MethodPatch, "https://www.google.com/search?t=test", nil)
	resp := &http.Response{
		Status:           "Not Found",
		StatusCode:       404,
		Proto:            "HTTP1.01",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	s := FormatLogJson("egress", start, time.Since(start), req, resp, "egress-route", 250, 100, 10, "95/1000", "true", "10%", "UT")

	fmt.Printf("test: FormatLogJson() -> %v\n", s)

	//Output:
	//fail

}

func _ExampleFormatLogText() {
	start := time.Now().UTC()
	req, _ := http.NewRequest(http.MethodPatch, "https://www.google.com/search?t=test", nil)
	resp := &http.Response{
		Status:           "Not Found",
		StatusCode:       404,
		Proto:            "HTTP1.01",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	s := FormatLogText("egress", start, time.Since(start), req, resp, "egress-route", 250, 100, 10, "95/1000", "true", "10%", "UT")

	fmt.Printf("test: FormatLogText() -> %v\n", s)

	//Output:
	//fail

}
