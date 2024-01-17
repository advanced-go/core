package runtime

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	searchResultsGzip = "file://[cwd]/runtimetest/search-results.gz"
)

func ExampleReadFile() {
	s := status504
	buf, status := ReadFile(s)
	fmt.Printf("test: ReadFile(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = address1Url
	buf, status = ReadFile(s)
	fmt.Printf("test: ReadFile(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(s), len(buf), status)

	s = status504
	u := parseRaw(s)
	buf, status = ReadFile(u.String())
	fmt.Printf("test: ReadFile(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	s = address1Url
	u = parseRaw(s)
	buf, status = ReadFile(u.String())
	fmt.Printf("test: ReadFile(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(u), len(buf), status)

	//Output:
	//test: ReadFile(file://[cwd]/runtimetest/status-504.json) -> [type:string] [buf:82] [status:OK]
	//test: ReadFile(file://[cwd]/runtimetest/address1.json) -> [type:string] [buf:68] [status:OK]
	//test: ReadFile(file://[cwd]/runtimetest/status-504.json) -> [type:*url.URL] [buf:82] [status:OK]
	//test: ReadFile(file://[cwd]/runtimetest/address1.json) -> [type:*url.URL] [buf:68] [status:OK]

}

func ExampleReadAll_Reader() {
	s := address3Url
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}

	r := strings.NewReader(string(buf0))
	buf, status := ReadAll(r, nil)
	fmt.Printf("test: ReadAll(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(r), len(buf), status)

	body := io.NopCloser(strings.NewReader(string(buf0)))
	buf, status = ReadAll(body, nil)
	fmt.Printf("test: ReadAll(%v) -> [type:%v] [buf:%v] [status:%v]\n", s, reflect.TypeOf(body), len(buf), status)

	//Output:
	//test: ReadAll(file://[cwd]/runtimetest/address3.json) -> [type:*strings.Reader] [buf:72] [status:OK]
	//test: ReadAll(file://[cwd]/runtimetest/address3.json) -> [type:io.nopCloserWriterTo] [buf:72] [status:OK]

}

func ExampleEncodingReader_Error() {
	s := address3Url
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	r := strings.NewReader(string(buf0))

	h := make(http.Header)
	h.Set(ContentEncoding, BrotliEncoding)
	reader, status := EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	h.Set(ContentEncoding, DeflateEncoding)
	reader, status = EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	h.Set(ContentEncoding, CompressEncoding)
	reader, status = EncodingReader(r, h)
	fmt.Printf("test: EncodingReader() -> [reader:%v] [status:%v]\n", reader, status)

	//Output:
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [br]]]
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [deflate]]]
	//test: EncodingReader() -> [reader:<nil>] [status:Content Decoding Failure [error: content encoding not supported [compress]]]

}

func ExampleReadAll_GzipReader() {
	s := searchResultsGzip
	buf0, err := os.ReadFile(FileName(s))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
	}
	var buf []byte
	status := StatusOK()

	direct := true
	h := make(http.Header)
	h.Set(ContentEncoding, GzipEncoding)
	r := bytes.NewReader(buf0)
	if direct {
		zr, err1 := EncodingReader(r, h)
		if !err1.OK() {
			fmt.Printf("gzip error: %v\n", err1)
			return
		}
		buf, err = io.ReadAll(zr)
		if err != nil {
			fmt.Printf("gzip error: %v\n", err)
			return
		}
		fmt.Printf("test: ReadAll() -> [buf:%v] [status:%v]\n", len(buf), status)
	} else {
		buf, status = ReadAll(r, h)
		fmt.Printf("test: ReadAll() -> [buf:%v] [status:%v]\n", len(buf), status)
		//buf, status = ReadAll(io.NopCloser(r), h)
		//fmt.Printf("test: ReadAll-NopCloser() -> [buf:%v] [status:%v]\n", len(buf), status) //s = string(buf)
	}

	//Output:
	//test: ReadAll() -> [buf:107600] [status:OK]

}

func ExampleReadAll_GzipReadCloser() {
	uri := "https://www.google.com/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)
	req.Header.Add(AcceptEncoding, "gzip, deflate, br")

	resp, err := http.DefaultClient.Do(req)
	fmt.Printf("test: Do() -> [content-type:%v] [content-encoding:%v] [err:%v]\n", resp.Header.Get(contentType), resp.Header.Get(ContentEncoding), err)

	buf, status := ReadAll(resp.Body, resp.Header)
	ct := http.DetectContentType(buf)
	fmt.Printf("test: ReadAll() -> [content-type:%v] [status:%v]\n", ct, status)

	//Output:
	//test: Do() -> [content-type:text/html; charset=ISO-8859-1] [content-encoding:gzip] [err:<nil>]
	//test: ReadAll() -> [content-type:text/html; charset=utf-8] [status:OK]

}
