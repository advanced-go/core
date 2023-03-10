package httptest

import (
	"fmt"
	"io"
	"strings"
)

func ExampleStringReader() {
	s := "This is an example of content"
	r := NewReaderCloser(strings.NewReader(s), nil)
	var buf = make([]byte, len(s))
	cnt, err := r.Read(buf)

	fmt.Printf("test: NewReaderCloser(s,nil) -> [error:%v] [cnt:%v] [content:%v]\n", err, cnt, string(buf))

	//Output:
	//test: NewReaderCloser(s,nil) -> [error:<nil>] [cnt:29] [content:This is an example of content]

}

func ExampleIOError() {
	s := "This is an example of content"
	r := NewReaderCloser(strings.NewReader(s), io.ErrUnexpectedEOF)
	var buf = make([]byte, len(s))
	cnt, err := r.Read(buf)

	fmt.Printf("test: NewReaderCloser(s,err) -> [error:%v] [cnt:%v] [content:%v]\n", err, cnt, len(buf) == 0)

	//Output:
	//test: NewReaderCloser(s,err) -> [error:unexpected EOF] [cnt:0] [content:false]

}
