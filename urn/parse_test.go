package urn

import (
	"fmt"
	"net/url"
)

func Example_Parse() {
	urn := "urn:github.com/advanced-go/example-domain/activity:entry"

	scheme, nsid, nss, err := Parse(urn)
	fmt.Printf("test: Parse(%v) ->[scheme:%v] [nsid:%v] [nss:%v] [err:%v]\n", urn, scheme, nsid, nss, err)

	/*
		u, err = url.Parse("github.com.advanced-go.example-domain.activity:/entry")
		fmt.Printf("test: url.Parse() -> [err:%v] [scheme:%v] [host:%v] [path:%v]\n", err, u.Scheme, u.Host, u.Path)

		u, err = url.Parse("go://github.com/advanced-go/example-domain/activity:entry")
		fmt.Printf("test: url.Parse() -> [err:%v] [scheme:%v] [host:%v] [path:%v]\n", err, u.Scheme, u.Host, u.Path)

		u, err = url.Parse("https://www.google.com/github.com/advanced-go/example-domain/activity:entry")
		fmt.Printf("test: url.Parse() -> [err:%v] [scheme:%v] [host:%v] [path:%v]\n", err, u.Scheme, u.Host, u.Path)

	*/

	//Output:
	//test: Parse(urn:github.com/advanced-go/example-domain/activity:entry) ->[scheme:urn] [nsid:github.com/advanced-go/example-domain/activity] [nss:entry] [err:<nil>]

}

func Example_Build() {
	urn := Build("test-nsid", "test-nss")

	fmt.Printf("test: Build() -> %v\n", urn)

	//Output:
	//test: Build() -> urn:test-nsid:test-nss

}

func Example_BuildUri() {
	urn := Build("github.com/advanced-go/example-domain/activity", "entry?q=golang")

	uri := ToUri("https", "www.google.com", urn)
	u, err := url.Parse(uri)
	fmt.Printf("test: ToUri() -> %v [url:%v] [err:%v]\n", urn, u.String(), err)

	urn1 := FromUri(u.String())
	fmt.Printf("test: FromUri() -> [url:%v] [urn:%v]\n", u.String(), urn1)

	//Output:
	//test: ToUri() -> urn:github.com/advanced-go/example-domain/activity:entry?q=golang [url:https://www.google.com/github.com/advanced-go/example-domain/activity:entry?q=golang] [err:<nil>]
	//test: FromUri() -> [url:https://www.google.com/github.com/advanced-go/example-domain/activity:entry?q=golang] [urn:urn:github.com/advanced-go/example-domain/activity:entry]

}
