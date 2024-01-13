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
)

func Example_Authority() {
	key := ""
	r := NewResolver() //WithAuthorities([]runtime.Pair{{"test", "www.bing.com"}})

	s, err := r.Authority(key)
	fmt.Printf("test: Authority-Empty(\"\") ->  [auth:%v] [err:%v]\n", s, err)

	key = "GOOGL"
	s, err = r.Authority(key)
	fmt.Printf("test: Authority-No-Variable-Markup(%v) ->  [auth:%v] [err:%v]\n", key, s, err)

	key = "{GOOGL"
	s, err = r.Authority(key)
	fmt.Printf("test: Authority-Invalid-Variable-Markup(%v) ->  [auth:%v] [err:%v]\n", key, s, err)

	key = "GOOGL}"
	s, err = r.Authority(key)
	fmt.Printf("test: Authority-Invalid-Variable-Markup(%v) ->  [auth:%v] [err:%v]\n", key, s, err)

	key = GOOGLVariable
	s, err = r.Authority(key)
	fmt.Printf("test: Authority-Valid-Variable(%v) ->  [auth:%v] [err:%v]\n", key, s, err)

	r.SetAuthorities([]runtime.Pair{{GOOGLVariable, GOOGLAuthority}, {MSFTVariable, MSFTAuthority}})
	key = GOOGLVariable
	s, err = r.Authority(key)
	fmt.Printf("test: Authority(%v) ->  [auth:%v] [err:%v]\n", key, s, err)

	//Output:
	//test: Authority-Empty("") ->  [auth:] [err:<nil>]
	//test: Authority-No-Variable-Markup(GOOGL) ->  [auth:GOOGL] [err:<nil>]
	//test: Authority-Invalid-Variable-Markup({GOOGL) ->  [auth:{GOOGL] [err:<nil>]
	//test: Authority-Invalid-Variable-Markup(GOOGL}) ->  [auth:GOOGL}] [err:<nil>]
	//test: Authority-Valid-Variable({GOOGL}) ->  [auth:] [err:resolver error: authority not found for variable: {GOOGL}]
	//test: Authority({GOOGL}) ->  [auth:www.google.com] [err:<nil>]

}

func Example_OverrideUrl() {
	r := NewResolverWithAuthorities([]runtime.Pair{{MSFTVariable, MSFTAuthority}})
	key := ""

	uri, ok := r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Empty(\"\") ->  [uri:%v] [ok:%v]\n", uri, ok)

	key = "www.google.com"
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Invalid-Variable(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	key = "{www.google.com"
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Invalid-Variable(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	key = "www.google.com}"
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Invalid-Variable(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	key = MSFTVariable
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Empty-Overrides(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	r.SetOverrides([]runtime.Pair{{MSFTVariable, "https://www.bing.com/office"}})
	key = MSFTVariable
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Configured-Overrrides(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	r.SetLocalHostOverride(true)
	key = MSFTVariable
	uri, ok = r.OverrideUrl(key)
	fmt.Printf("test: OverrideUrl-Configured-Overrides-LocalHost-Override(%v) ->  [uri:%v] [ok:%v]\n", key, uri, ok)

	//Output:
	//test: OverrideUrl-Empty("") ->  [uri:] [ok:false]
	//test: OverrideUrl-Invalid-Variable(www.google.com) ->  [uri:] [ok:false]
	//test: OverrideUrl-Invalid-Variable({www.google.com) ->  [uri:] [ok:false]
	//test: OverrideUrl-Invalid-Variable(www.google.com}) ->  [uri:] [ok:false]
	//test: OverrideUrl-Empty-Overrides({MSFT}) ->  [uri:] [ok:false]
	//test: OverrideUrl-Configured-Overrrides({MSFT}) ->  [uri:https://www.bing.com/office] [ok:true]
	//test: OverrideUrl-Configured-Overrides-LocalHost-Override({MSFT}) ->  [uri:https://www.bing.com/office] [ok:true]

}

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

func Example_Values() {
	v := make(url.Values)

	v.Add("param-1", "value-1")
	v.Add("param-2", "value-2")

	fmt.Printf("test: Values.Encode() -> %v\n", v.Encode())

	//Output:
	//test: Values.Encode() -> param-1=value-1&param-2=value-2

}
