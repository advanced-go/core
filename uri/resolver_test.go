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
	val := r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, val)

	id = "test"
	val = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, val)

	id = "/google/search?q=golang"
	val = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, val)

	id = "https://www.google.com/google:search?q=golang"
	val = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, val)

	//Output:
	//test: Resolve("") -> error: id cannot be resolved to a URL
	//test: Resolve("test") -> test
	//test: Resolve("/google/search?q=golang") -> http://localhost:8080/google/search?q=golang
	//test: Resolve("https://www.google.com/google:search?q=golang") -> https://www.google.com/google:search?q=golang

}

func Example_Resolver_Default() {
	r := NewResolver("http://localhost:8080", testDefault)

	v := make(url.Values)
	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	id := resolvedId
	url := r.Resolve(id, v)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = pathId
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = bypassId
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = googleUrl
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	//Output:
	//test: Resolve("resolved") -> http://localhost:8080/advanced-go/example-domain/activity:entry?param-1=value-1&param-2=value-2
	//test: Resolve("path") -> http://localhost:8080/advanced-go/example-domain/activity:entry
	//test: Resolve("bypass") -> bypass
	//test: Resolve("https://www.google.com/search?q=golang") -> https://www.google.com/search?q=golang

}

func Example_Resolver_Override() {
	r := NewResolver("http://localhost:8080", testDefault)
	r.SetOverride(testOverride, "")

	v := make(url.Values)
	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	id := resolvedId
	url := r.Resolve(id, v)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = pathId
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = bypassId
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	id = overrideBypassId
	url = r.Resolve(id, nil)
	fmt.Printf("test: Resolve(\"%v\") -> %v\n", id, url)

	//Output:
	//test: Resolve("resolved") -> file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt
	//test: Resolve("path") -> file://[cwd]/uri/uritest/html-response.txt
	//test: Resolve("bypass") -> bypass
	//test: Resolve("overrideBypass") -> http://localhost:8080/advanced-go/example-domain/activity:entry

}

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
