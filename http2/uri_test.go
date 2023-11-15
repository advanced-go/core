package http2

import (
	"fmt"
	"net/http"
)

/*
	func Example_ParsePkgUrl() {
		s := "github.com/advanced-go/core/runtime"
		u := ParsePkgUrl(s)
		fmt.Printf("test: ParsePkgUri(%v) %v\n", s, u)

		//Output:
		//test: ParsePkgUri(github.com/advanced-go/core/runtime) file://github.com/advanced-go/core/runtime

}
*/
func ExampleParseUri_Url() {
	uri := "www.google.com"
	scheme, host, path := ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "https://www.google.com"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "https://www.google.com/search?q=test"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	//Output:
	//test: ParseUri(www.google.com) -> [scheme:] [startup:] [path:www.google.com]
	//test: ParseUri(https://www.google.com) -> [scheme:https] [startup:www.google.com] [path:]
	//test: ParseUri(https://www.google.com/search?q=test) -> [scheme:https] [startup:www.google.com] [path:/search]

}

func ExampleParseUri_Urn() {
	uri := "urn"
	scheme, host, path := ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "urn:postgres"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	uri = "urn:postgres:query.access-log"
	scheme, host, path = ParseUri(uri)
	fmt.Printf("test: ParseUri(%v) -> [scheme:%v] [startup:%v] [path:%v]\n", uri, scheme, host, path)

	//Output:
	//test: ParseUri(urn) -> [scheme:] [startup:] [path:urn]
	//test: ParseUri(urn:postgres) -> [scheme:urn] [startup:postgres] [path:]
	//test: ParseUri(urn:postgres:query.access-log) -> [scheme:urn] [startup:postgres] [path:query.access-log]

}

func Example_BuildUrl() {
	template := "{scheme}://{host}{path}?{query}"
	uri := "http://localhost:8080/base-path/resource?first=false"
	req, _ := http.NewRequest("", uri, nil)

	url, err := BuildUrl(req.URL, template)
	fmt.Printf("test: OriginalUrl() -> %v\n", url)

	// Scheme
	template = "https://{host}{path}?{query}"
	url, err = BuildUrl(req.URL, template)
	fmt.Printf("test: BuildUrl(scheme) -> [error:%v] [%v]\n", err, url)

	// Exclude query
	template = "{scheme}://{host}{path}"
	url, err = BuildUrl(req.URL, template)
	fmt.Printf("test: BuildUrl(exclude-query) -> [error:%v] [%v]\n", err, url)

	// Host only
	template = "{scheme}://{host}"
	url, err = BuildUrl(req.URL, template)
	fmt.Printf("test: BuildUrl(startup-only) -> [error:%v] [%v]\n", err, url)

	// Scheme + startup
	template = "https://google.com{path}?{query}"
	url, err = BuildUrl(req.URL, template)
	fmt.Printf("test: BuildUrl(scheme+startup) -> [error:%v] [%v]\n", err, url)

	// Scheme + startup + path
	template = "https://google.com/search?{query}"
	url, err = BuildUrl(req.URL, template)
	fmt.Printf("test: BuildUrl(scheme+startup+path) -> [error:%v] [%v]\n", err, url)

	//Output:
	//test: OriginalUrl() -> http://localhost:8080/base-path/resource?first=false
	//test: BuildUrl(scheme) -> [error:<nil>] [https://localhost:8080/base-path/resource?first=false]
	//test: BuildUrl(exclude-query) -> [error:<nil>] [http://localhost:8080/base-path/resource]
	//test: BuildUrl(startup-only) -> [error:<nil>] [http://localhost:8080]
	//test: BuildUrl(scheme+startup) -> [error:<nil>] [https://google.com/base-path/resource?first=false]
	//test: BuildUrl(scheme+startup+path) -> [error:<nil>] [https://google.com/search?first=false]

}

func Example_BuildUrl_EmptyQuery() {
	template := "{scheme}://{host}{path}?{query}"
	uri := "http://localhost:8080/base-path/resource"
	req, _ := http.NewRequest("", uri, nil)

	url, _ := BuildUrl(req.URL, template)
	fmt.Printf("test: OriginalUrl() -> %v\n", url)

	//Output:
	//test: OriginalUrl() -> http://localhost:8080/base-path/resource

}
