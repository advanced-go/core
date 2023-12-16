package runtimetest

import (
	"context"
	"fmt"
	"net/http"
)

func Example_ContentLocation() {
	ctx := NewContentLocationContext(context.Background(), nil)
	uri, ok := ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h := make(http.Header)
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h.Add(ContentLocation, "https://www.google.com/search?q=golang")
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	h = make(http.Header)
	h.Add(ContentLocation, "file://[cwd]/runtimetest/test.txt")
	ctx = NewContentLocationContext(context.Background(), h)
	uri, ok = ContentLocationFromContext(ctx)
	fmt.Printf("test: ContentLocationFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:false] [uri:]
	//test: ContentLocationFromContext() ->  [ok:true] [uri:file://[cwd]/runtimetest/test.txt]

}

func Example_FileUrl_File() {
	ctx := NewFileUrlContext(context.Background(), "")
	uri, ok := FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url := "https://www.google.com/search?q=golang"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url = "file://[cwd]/runtimetest/test.txt"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:true] [uri:file://[cwd]/runtimetest/test.txt]

}

func Example_FileUrl_Urn() {
	ctx := NewFileUrlContext(context.Background(), "")
	uri, ok := FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url := "https://www.google.com/search?q=golang"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	url = "urn:status:ok"
	ctx = NewFileUrlContext(context.Background(), url)
	uri, ok = FileUrlFromContext(ctx)
	fmt.Printf("test: FileUrlFromContext() ->  [ok:%v] [uri:%v]\n", ok, uri)

	//Output:
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:false] [uri:]
	//test: FileUrlFromContext() ->  [ok:true] [uri:urn:status:ok]

}
