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

func Example_Authority() {
	r := NewResolver([]Attr{{"test", "www.microsoft.com"}})
	GOOGL := ""

	s, err := r.Authority(GOOGL)
	fmt.Printf("test: Authority(\"\") ->  [auth:%v] [err:%v]\n", s, err)

	GOOGL = "www.google.com"
	s, err = r.Authority(GOOGL)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", GOOGL, s, err)

	GOOGL = "{www.google.com"
	s, err = r.Authority(GOOGL)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", GOOGL, s, err)

	GOOGL = "www.google.com}"
	s, err = r.Authority(GOOGL)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", GOOGL, s, err)

	GOOGL = "{www.google.com}"
	s, err = r.Authority(GOOGL)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", GOOGL, s, err)

	GOOGL = "{test}"
	s, err = r.Authority(GOOGL)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", GOOGL, s, err)

	//Output:
	//test: Authority("") ->  [auth:] [err:<nil>]
	//test: Authority(www.google.com) ->  [auth:www.google.com] [err:<nil>]
	//test: Authority({www.google.com) ->  [auth:{www.google.com] [err:<nil>]
	//test: Authority(www.google.com}) ->  [auth:www.google.com}] [err:<nil>]
	//test: Authority({www.google.com}) ->  [auth:] [err:resolver error: authority not found for variable: {www.google.com}]
	//test: Authority({test}) ->  [auth:www.microsoft.com] [err:<nil>]

}

func Example_OverrideUrl() {
	r := NewResolver([]Attr{{"test", "www.microsoft.com"}})
	GOOGL := ""

	uri, ok := r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(\"\") ->  [uri:%v] [ok:%v]\n", uri, ok)

	GOOGL = "www.google.com"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	GOOGL = "{www.google.com"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	GOOGL = "www.google.com}"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	GOOGL = "{www.google.com}"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	GOOGL = "{test}"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	r.SetOverrides([]Attr{{"test", "https://www.microsoft.com/office"}})

	GOOGL = "{test}"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	r.SetLocalHostOverride(true)
	GOOGL = "{test}"
	uri, ok = r.OverrideUrl(GOOGL)
	fmt.Printf("test: OverrideUrl(%v) ->  [uri:%v] [ok:%v]\n", GOOGL, uri, ok)

	//Output:
	//test: OverrideUrl("") ->  [uri:] [ok:false]
	//test: OverrideUrl(www.google.com) ->  [uri:] [ok:false]
	//test: OverrideUrl({www.google.com) ->  [uri:] [ok:false]
	//test: OverrideUrl(www.google.com}) ->  [uri:] [ok:false]
	//test: OverrideUrl({www.google.com}) ->  [uri:] [ok:false]
	//test: OverrideUrl({test}) ->  [uri:] [ok:false]
	//test: OverrideUrl({test}) ->  [uri:https://www.microsoft.com/office] [ok:true]
	//test: OverrideUrl({test}) ->  [uri:https://www.microsoft.com/office] [ok:true]

}

func Example_Build_Error() {
	path := "/google/search"
	GOOGL := ""
	r := NewResolver(nil)

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

	r := NewResolver(nil)

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
	key := "MSFT"
	MSFT := fmt.Sprintf("{%v}", key)
	r := NewResolver([]Attr{{key, "www.microsoft.com"}})

	GOOGL := "{GOOGL}"
	uri := r.Build(GOOGL, path, "12345")
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", GOOGL, path, uri)

	uri = r.Build(MSFT, path, "12345")
	fmt.Printf("test: Build(\"%v\",\"%v\") -> [uri:%v]\n", MSFT, path, uri)

	r.SetLocalHostOverride(true)
	uri = r.Build(MSFT, path, "12345")
	fmt.Printf("test: Build-LocalHost(\"%v\",\"%v\") -> [uri:%v]\n", MSFT, path, uri)

	r.SetLocalHostOverride(false)
	r.SetOverrides([]Attr{{key, "https://www.microsoft.com/office/excel"}})
	uri = r.Build(MSFT, path, "12345")
	fmt.Printf("test: Build-Override(\"%v\",\"%v\") -> [uri:%v]\n", MSFT, path, uri)

	r.SetOverrides(nil)
	uri = r.Build(MSFT, path, "12345")
	fmt.Printf("test: Build-RemoveOverride(\"%v\",\"%v\") -> [uri:%v]\n", MSFT, path, uri)

	//Output:
	//test: Build("{GOOGL}","/some/resource/%v") -> [uri:resolver error: authority not found for variable: {GOOGL}]
	//test: Build("{MSFT}","/some/resource/%v") -> [uri:https://www.microsoft.com/some/resource/12345]
	//test: Build-LocalHost("{MSFT}","/some/resource/%v") -> [uri:http://localhost:8080/some/resource/12345]
	//test: Build-Override("{MSFT}","/some/resource/%v") -> [uri:https://www.microsoft.com/office/excel]
	//test: Build-RemoveOverride("{MSFT}","/some/resource/%v") -> [uri:https://www.microsoft.com/some/resource/12345]

}

/*

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


*/

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
