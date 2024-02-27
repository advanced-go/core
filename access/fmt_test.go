package access

import (
	"fmt"
	"net/http"
)

func ExampleSafeRequest() {
	r, _ := http.NewRequest("", "https://somehost.com/search?q=test", nil)
	r2 := SafeRequest(r)
	fmt.Printf("test: SafeRequest(r) -> [equal:%v]\n", r == r2)

	r = nil
	r2 = SafeRequest(r)
	fmt.Printf("test: SafeRequest(nil) -> [equal:%v]\n", r == r2)

	//Output:
	//test: SafeRequest(r) -> [equal:true]
	//test: SafeRequest(nil) -> [equal:false]

}

func ExampleSafeResponse() {
	r := new(http.Response)
	r2 := SafeResponse(r)
	fmt.Printf("test: SafeResponse(r) -> [equal:%v]\n", r == r2)

	r = nil
	r2 = SafeResponse(r)
	fmt.Printf("test: SafeResponse(nil) -> [equal:%v]\n", r == r2)

	//Output:
	//test: SafeResponse(r) -> [equal:true]
	//test: SafeResponse(nil) -> [equal:false]

}

func ExampleEncoding() {
	var r *http.Response
	e := Encoding(r)
	fmt.Printf("test: Encoding(nil) -> [encoding:%v]\n", e)

	r = new(http.Response)
	e = Encoding(r)
	fmt.Printf("test: Encoding(r) -> [encoding:%v]\n", e)

	r.Header = make(http.Header)
	r.Header.Add("Content-Encoding", "none")
	e = Encoding(r)
	fmt.Printf("test: Encoding(\"none\") -> [encoding:%v]\n", e)

	r.Header.Set("Content-Encoding", "gzip")
	e = Encoding(r)
	fmt.Printf("test: Encoding(\"gzip\") -> [encoding:%v]\n", e)

	//Output:
	//test: Encoding(nil) -> [encoding:]
	//test: Encoding(r) -> [encoding:]
	//test: Encoding("none") -> [encoding:]
	//test: Encoding("gzip") -> [encoding:gzip]

}
