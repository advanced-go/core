package exchange

import (
	"fmt"
	"net/http"
	"net/url"
)

func Example_Resolve() {
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

	resolver1 := Resolver{Valid: func(req *http.Request) bool {
		return req.URL.Path == uri1.Path
	}, Do: func(req *http.Request) (*http.Response, error) {
		return &resp1, nil
	}}

	resolver2 := Resolver{Valid: func(req *http.Request) bool {
		return req.URL.Path == uri2.Path
	}, Do: func(req *http.Request) (*http.Response, error) {
		return &resp2, nil
	}}

	err := AddResolver(resolver1)
	fmt.Printf("test: AddResolver(resolver1) -> [err:%v]\n", err)

	err = AddResolver(resolver2)
	fmt.Printf("test: AddResolver(resolver2) -> [err:%v]\n", err)

	req, _ := http.NewRequest("", "http://localhost:8080/endpoint/invalid", nil)
	do := Resolve(req)
	fmt.Printf("test: Resolve(%v) [do:%v]\n", req.URL.Path, do != nil)

	var resp *http.Response

	req, _ = http.NewRequest("", uri1.String(), nil)
	do = Resolve(req)
	if do != nil {
		resp, _ = do(req)
	}
	fmt.Printf("test: Resolve(%v) [do:%v] [statusCode:%v]\n", req.URL.Path, do != nil, resp.StatusCode)

	req, _ = http.NewRequest("", uri2.String(), nil)
	do = Resolve(req)
	if do != nil {
		resp, _ = do(req)
	}
	fmt.Printf("test: Resolve(%v) [do:%v] [statusCode:%v]\n", req.URL.Path, do != nil, resp.StatusCode)

	//Output:
	//test: AddResolver(resolver1) -> [err:<nil>]
	//test: AddResolver(resolver2) -> [err:<nil>]
	//test: Resolve(/endpoint/invalid) [do:false]
	//test: Resolve(/endpoint1) [do:true] [statusCode:200]
	//test: Resolve(/endpoint2) [do:true] [statusCode:404]

}
