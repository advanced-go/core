package http2

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

func Example_WriteBytes_Error() {
	var cnt = 100
	var pstr *string
	var pr *io.Reader

	buf, rc, status := WriteBytes(nil, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, rc)

	buf, rc, status = WriteBytes(pstr, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, rc)

	buf, rc, status = WriteBytes(pr, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, rc)

	buf, rc, status = WriteBytes(cnt, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, rc)

	str := "this error should not be seen"
	buf, rc, status = WriteBytes(NewReaderCloser(strings.NewReader(str), errors.New("error on internal read")), "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [rc:%v]\n", buf != nil, status, rc)

	//Output:
	//test: WriteBytes() -> [buf:false] [status:Internal Error [error: content type is invalid: <nil>]] [rc:]
	//test: WriteBytes() -> [buf:false] [status:Internal Error [error: content type is invalid: *string]] [rc:]
	//test: WriteBytes() -> [buf:false] [status:Internal Error [error: content type is invalid: *io.Reader]] [rc:]
	//test: WriteBytes() -> [buf:false] [status:Internal Error [error: content type is invalid: int]] [rc:]
	//test: WriteBytes() -> [buf:false] [status:I/O Failure [error on internal read]] [rc:]

}

func Example_WriteBytes() {
	str := "this is string content"
	buf, rc, status := WriteBytes(str, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	str = "this is []byte content"
	buf, rc, status = WriteBytes([]byte(str), "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	str = "this is an error message"
	buf, rc, status = WriteBytes(errors.New(str), "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	str = "this is http.Response.Body content"
	r := strings.NewReader(str)
	resp := new(http.Response)
	resp.Body = io.NopCloser(r)
	buf, rc, status = WriteBytes(resp, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	str = "this is io.Reader content"
	r = strings.NewReader(str)
	buf, rc, status = WriteBytes(r, "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	str = "this is io.ReaderCloser content"
	r = strings.NewReader(str)
	buf, rc, status = WriteBytes(io.NopCloser(r), "")
	fmt.Printf("test: WriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	//Output:
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is string content] [rc:text/plain; charset=utf-8]
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is []byte content] [rc:text/plain; charset=utf-8]
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is an error message] [rc:text/plain; charset=utf-8]
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is http.Response.Body content] [rc:text/plain; charset=utf-8]
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is io.Reader content] [rc:text/plain; charset=utf-8]
	//test: WriteBytes() -> [buf:true] [status:OK] [content:this is io.ReaderCloser content] [rc:text/plain; charset=utf-8]

}

func Example_WriteBytes_Json() {
	c := data{"123456", 101}
	buf, rc, status := WriteBytes(c, ContentTypeJson)
	fmt.Printf("test: WriteWriteBytes() -> [buf:%v] [status:%v] [content:%v] [rc:%v]\n", buf != nil, status, string(buf), rc)

	//Output:
	//test: WriteWriteBytes() -> [buf:true] [status:OK] [content:{"Item":"123456","Count":101}] [rc:application/json]

}
