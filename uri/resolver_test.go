package uri

import (
	"fmt"
	"net/url"
)

const (
	resolvedId       = "resolved"
	bypassId         = "bypass"
	overrideBypassId = "overrideBypass"
	pathId           = "path"

	activityUrl  = "http://localhost:8080/advanced-go/example-domain/activity:entry"
	activityPath = "/advanced-go/example-domain/activity:entry"
	googleUrl    = "https://www.google.com/search?q=golang"
	//googlePath   = "/serach?q=golang"

	fileUrl  = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	filePath = "file://[cwd]/uri/uritest/html-response.txt"
)

func testDefault(id string) string {
	switch id {
	case resolvedId:
		return activityUrl
	case pathId:
		return activityPath
	case bypassId:
		return ""
	case overrideBypassId:
		return activityPath
	}
	return id
}

func testOverride(id string) (string, string) {
	switch id {
	case resolvedId:
		return fileUrl, ""
	case pathId:
		return filePath, ""
	case overrideBypassId:
		return "", ""
	}
	return id, ""
}

func Example_Resolver_Passthrough() {
	r := NewResolver("http://localhost:8080", nil)

	id := ""
	val := r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, val)

	id = "test"
	val = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, val)

	id = "/google/search?q=golang"
	val = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, val)

	id = "https://www.google.com/google:search?q=golang"
	val = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, val)

	//Output:
	//test: Build("") -> error: id cannot be resolved to a URL
	//test: Build("test") -> test
	//test: Build("/google/search?q=golang") -> http://localhost:8080/google/search?q=golang
	//test: Build("https://www.google.com/google:search?q=golang") -> https://www.google.com/google:search?q=golang

}

func Example_Resolver_Default() {
	r := NewResolver("http://localhost:8080", testDefault)

	v := make(url.Values)
	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	id := resolvedId
	url := r.Build(id, v)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = pathId
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = bypassId
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = googleUrl
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	//Output:
	//test: Build("resolved") -> http://localhost:8080/advanced-go/example-domain/activity:entry?param-1=value-1&param-2=value-2
	//test: Build("path") -> http://localhost:8080/advanced-go/example-domain/activity:entry
	//test: Build("bypass") -> bypass
	//test: Build("https://www.google.com/search?q=golang") -> https://www.google.com/search?q=golang

}

func Example_Resolver_Override() {
	r := NewResolver("http://localhost:8080", testDefault)
	r.SetOverride(testOverride, "")

	v := make(url.Values)
	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	id := resolvedId
	url := r.Build(id, v)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = pathId
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = bypassId
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	id = overrideBypassId
	url = r.Build(id, nil)
	fmt.Printf("test: Build(\"%v\") -> %v\n", id, url)

	//Output:
	//test: Build("resolved") -> file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt
	//test: Build("path") -> file://[cwd]/uri/uritest/html-response.txt
	//test: Build("bypass") -> bypass
	//test: Build("overrideBypass") -> http://localhost:8080/advanced-go/example-domain/activity:entry

}

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
