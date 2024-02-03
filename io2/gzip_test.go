package io2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	testResponseHtml = "file://[cwd]/io2test/test-response.txt"
	testResponseGzip = "file://[cwd]/io2test/test-response.gz"
)

func ExampleGzipReader() {
	content, err0 := os.ReadFile(FileName(testResponseGzip))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}

	// read content
	zr, status := NewGzipReader(bytes.NewReader(content))
	fmt.Printf("test: NewGzipReader() -> [status:%v]\n", status)

	buff, err1 := io.ReadAll(zr)
	status = zr.Close()
	fmt.Printf("test: ReadAll(gzip.Reader()) -> [read-err:%v] [close-status:%v]\n", err1, status)
	fmt.Printf("test: DetectContent -> [input:%v] [output:%v] [out-len:%v]\n", http.DetectContentType(content), http.DetectContentType(buff), len(buff))

	//Output:
	//test: NewGzipReader() -> [status:OK]
	//test: ReadAll(gzip.Reader()) -> [read-err:<nil>] [close-status:OK]
	//test: DetectContent -> [input:application/x-gzip] [output:text/plain; charset=utf-8] [out-len:188]

}

func _ExampleGzipWriter() {
	content, err0 := os.ReadFile(FileName(testResponseHtml))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buff := new(bytes.Buffer)

	// write content
	zw := NewGzipWriter(buff)
	cnt, err := zw.Write(content)
	status := zw.Close()
	fmt.Printf("test: gzip.Writer() -> [cnt:%v] [write-err:%v] [close-status:%v]\n", cnt, err, status)

	buff2 := bytes.Clone(buff.Bytes())
	fmt.Printf("test: DetectContent -> [input:%v] [output:%v]\n", http.DetectContentType(content), http.DetectContentType(buff2))

	err = os.WriteFile(FileName(testResponseGzip), buff2, 667)
	fmt.Printf("test: os.WriteFile(\"%v\") -> [err:%v]\n", testResponseGzip, err)

	//Output:
	//test: gzip.Writer() -> [cnt:188] [write-err:<nil>] [close-status:OK]
	//test: DetectContent -> [input:text/plain; charset=utf-8] [output:application/x-gzip]
	//test: os.WriteFile("file://[cwd]/runtimetest/test-response.gz") -> [err:<nil>]

}
