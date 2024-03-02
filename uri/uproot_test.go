package uri

import (
	"fmt"
	"net/url"
)

func ExampleParseRaw() {
	u := "http://localhost:8080/github/advanced-go/core/uri.Uproot?q=golang"
	uri, err := url.Parse(u)

	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [host:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Host, uri.Path, uri.Fragment, uri.RawQuery, err)

	u = "http://localhost:8080/github/advanced-go/core/uri:Uproot?q=golang"
	uri, err = url.Parse(u)

	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Path, uri.Fragment, uri.RawQuery, err)

	u = "http://localhost:8080/github/advanced-go/core/uri?q=golang#Uproot"
	uri, err = url.Parse(u)

	fmt.Printf("test: ParseRaw(\"%v\") -> [scheme:%v] [path:%v] [frag:%v] [query:%v] [err:%v]\n", u, uri.Scheme, uri.Path, uri.Fragment, uri.RawQuery, err)

	//Output:
}

func Example_Uproot() {
	uri := ""
	nid, nss, ok := UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	// Urn should not be changed
	uri = "urn:github.resource"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	// valid path only and an empty nss
	uri = "/valid-empty-nss?q=golang"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	// valid embedded path only
	uri = "/github/valid-leading-slash/example-domain/activity:entry"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	// valid URN
	uri = "github.com/valid-no-leading-slash/example-domain/activity:entry"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	uri = "https://www.google.com/valid-uri?q=golang"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	uri = "https://www.google.com/github.com/valid-uri-nss/search?q=golang"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	uri = "https://www.google.com/github.com/valid-uri-with-nss:search?q=golang"
	nid, nss, ok = UprootUrn(uri)
	fmt.Printf("test: Uproot(%v) -> [nid:%v] [nss:%v] [ok:%v]\n", uri, nid, nss, ok)

	//Output:
	//test: Uproot() -> [nid:] [nss:] [ok:false]
	//test: Uproot(urn:github.resource) -> [nid:urn:github.resource] [nss:] [ok:true]
	//test: Uproot(/valid-empty-nss?q=golang) -> [nid:valid-empty-nss] [nss:] [ok:true]
	//test: Uproot(/github/valid-leading-slash/example-domain/activity:entry) -> [nid:github/valid-leading-slash/example-domain/activity] [nss:entry] [ok:true]
	//test: Uproot(github.com/valid-no-leading-slash/example-domain/activity:entry) -> [nid:github.com/valid-no-leading-slash/example-domain/activity] [nss:entry] [ok:true]
	//test: Uproot(https://www.google.com/valid-uri?q=golang) -> [nid:valid-uri] [nss:] [ok:true]
	//test: Uproot(https://www.google.com/github.com/valid-uri-nss/search?q=golang) -> [nid:github.com/valid-uri-nss/search] [nss:] [ok:true]
	//test: Uproot(https://www.google.com/github.com/valid-uri-with-nss:search?q=golang) -> [nid:github.com/valid-uri-with-nss] [nss:search] [ok:true]
}
