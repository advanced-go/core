package io2

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	testResponseTxt  = "file://[cwd]/io2test/test-response.txt"
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
	err2 := zr.Close()
	fmt.Printf("test: ReadAll(gzip.Reader()) -> [read-err:%v] [close-err:%v]\n", err1, err2)
	fmt.Printf("test: DetectContent -> [input:%v] [output:%v] [out-len:%v]\n", http.DetectContentType(content), http.DetectContentType(buff), len(buff))

	//Output:
	//test: NewGzipReader() -> [status:OK]
	//test: ReadAll(gzip.Reader()) -> [read-err:<nil>] [close-err:<nil>]
	//test: DetectContent -> [input:application/x-gzip] [output:text/plain; charset=utf-8] [out-len:188]

}

func ExampleGzipWriter() {
	content, err0 := os.ReadFile(FileName(testResponseTxt))
	if err0 != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err0)
		return
	}
	buff := new(bytes.Buffer)

	// write content
	zw := NewGzipWriter(buff)
	cnt, err := zw.Write(content)
	err1 := zw.Close()
	fmt.Printf("test: gzip.Writer() -> [cnt:%v] [write-err:%v] [close-err:%v]\n", cnt, err, err1)

	buff2 := bytes.Clone(buff.Bytes())
	fmt.Printf("test: DetectContent -> [input:%v] [output:%v]\n", http.DetectContentType(content), http.DetectContentType(buff2))

	err = os.WriteFile(FileName(testResponseGzip), buff2, 667)
	fmt.Printf("test: os.WriteFile(\"%v\") -> [err:%v]\n", testResponseGzip, err)

	//Output:
	//test: gzip.Writer() -> [cnt:188] [write-err:<nil>] [close-err:<nil>]
	//test: DetectContent -> [input:text/plain; charset=utf-8] [output:application/x-gzip]
	//test: os.WriteFile("file://[cwd]/io2test/test-response.gz") -> [err:<nil>]

}
