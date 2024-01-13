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

	fileUrl   = "file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt"
	fileAttrs = "file://[cwd]/uritest/attrs.json"

	yahooSearch = "https://search.yahoo.com/search?p=golang"
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

func Example_Build() {
	path := ""
	r := NewResolver()

	uri := r.Build(path)
	fmt.Printf("test: Build(\"%v\") -> [uri:%v]\n", path, uri)

	path = "/search?q=golang"
	uri = r.Build(path)
	fmt.Printf("test: Build(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides([]runtime.Pair{{path, yahooSearch}})
	uri = r.Build(path)
	fmt.Printf("test: Build(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides(nil)
	values := make(url.Values)
	values.Add("q", "golang")
	path = "/search?%v"
	uri = r.Build(path, values.Encode())
	fmt.Printf("test: Build(\"%v\") -> [uri:%v]\n", path, uri)

	r.SetOverrides([]runtime.Pair{{path, yahooSearch}})
	uri = r.Build(path, values.Encode())
	fmt.Printf("test: Build(\"%v\") -> [uri:%v]\n", path, uri)

	//Output:
	//test: Build("") -> [uri:resolver error: invalid argument, path is empty]
	//test: Build("/search?q=golang") -> [uri:http://localhost:8080/search?q=golang]
	//test: Build("/search?q=golang") -> [uri:https://search.yahoo.com/search?p=golang]
	//test: Build("/search?%v") -> [uri:http://localhost:8080/search?q=golang]
	//test: Build("/search?%v") -> [uri:https://search.yahoo.com/search?p=golang]

}

/*
func Example_Build_Error() {
	path := "/google/search"
	GOOGL := ""
	r := NewResolver()

	uri := r.Build(GOOGL, path)
	fmt.Printf("test: Build(\"\",\"%v\") -> [uri:%v]\n", path, uri)

	GOOGL = "www.google.com"
	uri = r.Build(GOOGL, "")
	fmt.Printf("test: Build(\"%v\",\"\") -> [uri:%v]\n", GOOGL, uri)

	//Output:
	//test: Build("","/google/search") -> [uri:resolver error: invalid argument, authority is empty]
	//test: Build("www.google.com","") -> [uri:resolver error: invalid argument, path is empty]

}

func Example_Build_Passthrough() {
	GOOGL := "www.google.com"
	path := "/%v?%v"
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	r := NewResolver()

	enc := v.Encode()
	uri := r.Build(GOOGL, path, "search", enc)
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", GOOGL, path, uri)

	r.SetLocalHostOverride(true)
	uri = r.Build(GOOGL, path, "search", enc)
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", GOOGL, path, uri)

	//Output:
	//test: Build("www.google.com","/%v?%v") -> [uri:https://www.google.com/search?param-1=value-1&param-2=value-2]
	//test: Build("www.google.com","/%v?%v") -> [uri:http://localhost:8080/search?param-1=value-1&param-2=value-2]

}

func Example_Build() {
	path := "/some/resource/%v"
	r := NewResolverWithAuthorities([]runtime.Pair{{MSFTVariable, MSFTAuthority}})

	key := GOOGLVariable
	uri := r.Build(key, path, "12345")
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	key = MSFTVariable
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	r.SetAuthorities([]runtime.Pair{{GOOGLVariable, GOOGLAuthority}})
	key = MSFTVariable
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build-SetAuthorities(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	r.SetAuthorities([]runtime.Pair{{MSFTVariable, MSFTAuthority}})
	key = MSFTVariable
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build-ResetAuthorities(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	r.SetLocalHostOverride(true)
	key = MSFTVariable
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build-LocalHost(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	r.SetLocalHostOverride(false)
	key = MSFTVariable
	r.SetOverrides([]runtime.Pair{{MSFTVariable, fileUrl}})
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build-Override(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	r.SetOverrides([]runtime.Pair{{GOOGLVariable, GOOGLAuthority}})
	uri = r.Build(key, path, "12345")
	fmt.Printf("test: Build-RemoveOverride(\"%v\",\"%v\") -> [uri:%v]\n", key, path, uri)

	//Output:
	//test: Build("{GOOGL}","/some/resource/%v") -> [uri:resolver error: authority not found for variable: {GOOGL}]
	//test: Build("{MSFT}","/some/resource/%v") -> [uri:https://www.bing.com/some/resource/12345]
	//test: Build-SetAuthorities("{MSFT}","/some/resource/%v") -> [uri:resolver error: authority not found for variable: {MSFT}]
	//test: Build-ResetAuthorities("{MSFT}","/some/resource/%v") -> [uri:https://www.bing.com/some/resource/12345]
	//test: Build-LocalHost("{MSFT}","/some/resource/%v") -> [uri:http://localhost:8080/some/resource/12345]
	//test: Build-Override("{MSFT}","/some/resource/%v") -> [uri:file:///c:/Users/markb/GitHub/core/uri/uritest/html-response.txt]
	//test: Build-RemoveOverride("{MSFT}","/some/resource/%v") -> [uri:https://www.bing.com/some/resource/12345]

}


*/

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
