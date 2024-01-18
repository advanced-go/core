package runtime

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type data struct {
	Item  string
	Count int
}

// ReaderCloser - test type for a body.ReadCloser interface
type ReaderCloser struct {
	Reader io.Reader
	Err    error
}

func (r *ReaderCloser) Read(p []byte) (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	return r.Reader.Read(p)
}

func (r *ReaderCloser) Close() error {
	return nil
}

func NewReaderCloser(reader io.Reader, err error) *ReaderCloser {
	rc := new(ReaderCloser)
	rc.Reader = reader
	rc.Err = err
	return rc
}

func ExampleBytes_Error() {
	var cnt = 100
	var pstr *string
	var pr *io.Reader

	buf, status := Bytes(nil, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, http.DetectContentType(buf))

	buf, status = Bytes(pstr, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, http.DetectContentType(buf))

	buf, status = Bytes(pr, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, http.DetectContentType(buf))

	buf, status = Bytes(cnt, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, http.DetectContentType(buf))

	str := "this error should not be seen"
	buf, status = Bytes(NewReaderCloser(strings.NewReader(str), errors.New("error on internal read")), "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, http.DetectContentType(buf))

	//Output:
	//test: Bytes() -> [buf:false] [status:Internal Error [error: content type is invalid: <nil>]] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:false] [status:Internal Error [error: content type is invalid: *string]] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:false] [status:Internal Error [error: content type is invalid: *io.Reader]] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:false] [status:Internal Error [error: content type is invalid: int]] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:false] [status:I/O Failure [error on internal read]] [rc:text/plain; charset=utf-8]

}

func ExampleBytes() {
	str := "this is string content"
	buf, status := Bytes(str, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	str = "this is []byte content"
	buf, status = Bytes([]byte(str), "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	str = "this is an error message"
	buf, status = Bytes(errors.New(str), "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	//str = "this is http.Response.Body content"
	//r := strings.NewReader(str)
	//resp := new(http.Response)
	//resp.Body = io.NopCloser(r)
	//buf, status = Bytes(resp, "")
	//fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	str = "this is io.Reader content"
	r := strings.NewReader(str)
	buf, status = Bytes(r, "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	str = "this is io.ReaderCloser content"
	r = strings.NewReader(str)
	buf, status = Bytes(io.NopCloser(r), "")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), http.DetectContentType(buf))

	//Output:
	//test: Bytes() -> [buf:true] [status:OK] [content:this is string content] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:true] [status:OK] [content:this is []byte content] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:true] [status:OK] [content:this is an error message] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:true] [status:OK] [content:this is io.Reader content] [rc:text/plain; charset=utf-8]
	//test: Bytes() -> [buf:true] [status:OK] [content:this is io.ReaderCloser content] [rc:text/plain; charset=utf-8]

}

func ExampleBytes_Json() {
	c := data{"123456", 101}
	buf, status := Bytes(c, "application/json")
	fmt.Printf("test: Bytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), "application/json")

	//Output:
	//test: Bytes() -> [buf:true] [status:OK] [content:{"Item":"123456","Count":101}] [rc:application/json]

}
