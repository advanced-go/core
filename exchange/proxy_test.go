package exchange

import (
	"fmt"
	"net/http"
	"net/url"
)

func Example_ProxyLookup() {
	uri1, _ := url.Parse("http://localhost:8080/endpoint1")
	resp1 := http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		Proto:            "",
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
	uri2, _ := url.Parse("http://localhost:8080/endpoint2")
	resp2 := http.Response{
		Status:           "Not Found",
		StatusCode:       http.StatusNotFound,
		Proto:            "",
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

	proxy1 := Proxy{Select: func(req *http.Request) bool {
		return req.URL.Path == uri1.Path
	}, Do: func(req *http.Request) (*http.Response, error) {
		return &resp1, nil
	}}

	proxy2 := Proxy{Select: func(req *http.Request) bool {
		return req.URL.Path == uri2.Path
	}, Do: func(req *http.Request) (*http.Response, error) {
		return &resp2, nil
	}}

	err := AddProxy(proxy1)
	fmt.Printf("test: AddProxy(proxy1) -> [err:%v]\n", err)

	err = AddProxy(proxy2)
	fmt.Printf("test: AddProxy(proxy2) -> [err:%v]\n", err)

	req, _ := http.NewRequest("", "http://localhost:8080/endpoint/invalid", nil)
	do := ProxyLookup(req)
	fmt.Printf("test: ProxyLookup(%v) [do:%v]\n", req.URL.Path, do != nil)

	var resp *http.Response

	req, _ = http.NewRequest("", uri1.String(), nil)
	do = ProxyLookup(req)
	if do != nil {
		resp, _ = do(req)
	}
	fmt.Printf("test: ProxyLookup(%v) [do:%v] [statusCode:%v]\n", req.URL.Path, do != nil, resp.StatusCode)

	req, _ = http.NewRequest("", uri2.String(), nil)
	do = ProxyLookup(req)
	if do != nil {
		resp, _ = do(req)
	}
	fmt.Printf("test: ProxyLookup(%v) [do:%v] [statusCode:%v]\n", req.URL.Path, do != nil, resp.StatusCode)

	//Output:
	//test: AddProxy(proxy1) -> [err:<nil>]
	//test: AddProxy(proxy2) -> [err:<nil>]
	//test: ProxyLookup(/endpoint/invalid) [do:false]
	//test: ProxyLookup(/endpoint1) [do:true] [statusCode:200]
	//test: ProxyLookup(/endpoint2) [do:true] [statusCode:404]

}
