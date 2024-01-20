package runtime

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
)

const (
	searchResultsGzip = "file://[cwd]/runtimetest/search-results.gz"
	testResponseHtml  = "file://[cwd]/runtimetest/test-response.html"
	testResponseGzip  = "file://[cwd]/runtimetest/test-response.gz"
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

/*
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
*/
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

func _ExampleGzipWriter() {
	content, err0 := os.ReadFile(FileName(testResponseHtml))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buff := new(bytes.Buffer)

	// write, flush and close
	zw := gzip.NewWriter(buff)
	cnt, err := zw.Write(content)
	ferr := zw.Flush()
	cerr := zw.Close()
	fmt.Printf("test: gzip.Writer() -> [cnt:%v] [write-err:%v] [flush-err:%v] [close_err:%v]\n", cnt, err, ferr, cerr)

	buff2 := bytes.Clone(buff.Bytes())
	fmt.Printf("test: DetectContent -> [input:%v] [output:%v]\n", http.DetectContentType(content), http.DetectContentType(buff2))

	err = os.WriteFile(FileName(testResponseGzip), buff2, 667)
	fmt.Printf("test: os.WriteFile(\"%v\") -> [err:%v]\n", testResponseGzip, err)

	// decode the content
	//r := bytes.NewReader(buff2)
	//zr, rerr := gzip.NewReader(r)
	//buff1, err1 := io.ReadAll(zr)
	//cerr = zr.Close()
	//fmt.Printf("test: gzip.Reader() -> [new-err:%v] [read-err:%v] [close-err:%v]\n", rerr, err1, cerr)

	//fmt.Printf("test: DetectContent -> [input:%v] [output:%v] [out-len:%v]\n", http.DetectContentType(buff2), http.DetectContentType(buff1), len(buff1))

	//Output:
	//test: gzip.Writer() -> [cnt:188] [write-err:<nil>] [flush-err:<nil>] [close_err:<nil>]
	//test: DetectContent -> [input:text/plain; charset=utf-8] [output:application/x-gzip]
	//test: os.WriteFile("file://[cwd]/runtimetest/test-response.gz") -> [err:<nil>]

}

func ExampleGzipReader() {
	content, err0 := os.ReadFile(FileName(testResponseGzip))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	// decode the content
	r := bytes.NewReader(content)
	zr, err := gzip.NewReader(r)
	buff, err1 := io.ReadAll(zr)
	cerr := zr.Close()
	fmt.Printf("test: gzip.Reader() -> [new-err:%v] [read-err:%v] [close-err:%v]\n", err, err1, cerr)

	fmt.Printf("test: DetectContent -> [input:%v] [output:%v] [out-len:%v]\n", http.DetectContentType(content), http.DetectContentType(buff), len(buff))

	//Output:
	//test: gzip.Reader() -> [new-err:<nil>] [read-err:<nil>] [close-err:<nil>]
	//test: DetectContent -> [input:application/x-gzip] [output:text/plain; charset=utf-8] [out-len:188]

}
