package exchange

import (
	"fmt"
	"strings"
)

func ExampleStringReader() {
	s := "This is an example of content"
	r := NewReaderCloser(strings.NewReader(s))
	var buf = make([]byte, len(s))
	cnt, err := r.Read(buf)

	fmt.Printf("test: NewReaderCloser(s,nil) -> [error:%v] [cnt:%v] [content:%v]\n", err, cnt, string(buf))

	//Output:
	//test: NewReaderCloser(s,nil) -> [error:<nil>] [cnt:29] [content:This is an example of content]

}
