package uri

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

const (
	MSFTVariable  = "{MSFT}"
	MSFTAuthority = "www.bing.com"

	GOOGLVariable  = "{GOOGL}"
	GOOGLAuthority = "www.google.com"

	fileUrl   = "file:///c:/Users/markb/GitHub/core/uri/uritest/test-response.html"
	fileAttrs = "file://[cwd]/uritest/attrs.json"

	yahooSearch         = "https://search.yahoo.com/search?p=golang"
	yahooSearchTemplate = "https://search.yahoo.com/search?%v"
)

func Example_OverrideUrl() {
	r := NewResolver()
	path := ""

	uri, ok := r.OverrideUrl(path)
	fmt.Printf("test: OverrideUrl-Empty(\"\") ->  [uri:%v] [ok:%v]\n", uri, ok)

	path = "/search"
	uri, ok = r.OverrideUrl(path)
	fmt.Printf("test: OverrideUrl-Invalid-Path(\"%v\") ->  [uri:%v] [ok:%v]\n", path, uri, ok)

	path = "/search"
	r.SetOverrides([]runtime.Pair{{path, yahooSearch}})
	uri, ok = r.OverrideUrl(path)
	fmt.Printf("test: OverrideUrl-Valid(\"%v\") ->  [uri:%v] [ok:%v]\n", path, uri, ok)

	//Output:
	//test: OverrideUrl-Empty("") ->  [uri:] [ok:false]
	//test: OverrideUrl-Invalid-Path("/search") ->  [uri:] [ok:false]
	//test: OverrideUrl-Valid("/search") ->  [uri:https://search.yahoo.com/search?p=golang] [ok:true]

}

func ExampleBuild() {
	path := ""
	r := NewResolver()

	uri := r.Build(path)
	fmt.Printf("test: Build-Error(\"%v\") -> [uri:%v]\n", path, uri)

	path = "/search?q=golang"
	uri = r.Build(path)
	fmt.Printf("test: Build-Default(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides([]runtime.Pair{{path, yahooSearch}})
	uri = r.Build(path)
	fmt.Printf("test: Build-Override(\"%v\") -> [uri:%v]\n", path, uri)

	//Output:
	//test: Build-Error("") -> [uri:resolver error: invalid argument, path is empty]
	//test: Build-Default("/search?q=golang") -> [uri:http://localhost:8080/search?q=golang]
	//test: Build-Override("/search?q=golang") -> [uri:https://search.yahoo.com/search?p=golang]

}

func ExampleBuild_Values() {
	path := ""
	r := NewResolver()

	values := make(url.Values)
	values.Add("q", "golang")
	path = "/search?%v"
	uri := r.Build(path, values.Encode())
	fmt.Printf("test: Build-Values(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides([]runtime.Pair{{path, yahooSearchTemplate}})
	uri = r.Build(path, values.Encode())
	fmt.Printf("test: Build-Override-Values(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides([]runtime.Pair{{path, fileAttrs}})
	uri = r.Build(path, values.Encode())
	fmt.Printf("test: Build-Override-File-Scheme(\"%v\") -> [uri:%v]\n", path, uri)

	//Output:
	//test: Build-Values("/search?%v") -> [uri:http://localhost:8080/search?q=golang]
	//test: Build-Override-Values("/search?%v") -> [uri:https://search.yahoo.com/search?q=golang]
	//test: Build-Override-File-Scheme("/search?%v") -> [uri:file://[cwd]/uritest/attrs.json]

}

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
