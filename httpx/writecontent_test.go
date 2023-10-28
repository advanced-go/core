package httpx

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/runtime/runtimetest"
	"io"
	"strings"
)

type data struct {
	Item  string
	Count int
}

func Example_writeStatusContent() {
	r := NewRecorder()

	// No content
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatusOK(), "test location")
	r.Result().Header = r.Header()
	buf, status := ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Error message
	r = NewRecorder()
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatus(runtime.StatusInternal).SetContent("error message"), "test location")
	r.Result().Header = r.Header()
	buf, status = ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Json
	c := data{"123456", 101}
	r = NewRecorder()
	writeStatusContent[runtimetest.DebugError](r, runtime.NewStatus(runtime.StatusInternal).SetContent(c), "test location")
	r.Result().Header = r.Header()
	buf, status = ReadAll(r.Result().Body)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	//Output:
	//test: writeStatusContent() -> 200 [header:map[]] [body:] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[13] Content-Type:[text/plain; charset=utf-8]]] [body:error message] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Length:[29] Content-Type:[application/json]]] [body:{"Item":"123456","Count":101}] [ReadAll:OK]
}

func Example_serializeContent() {
	var cnt = 100

	buf, status := serializeContent[int](cnt)
	fmt.Printf("test: serializeContent() -> [buf:%v] [status:%v]\n", buf != nil, status)

	str := "this is string content"
	buf, status = serializeContent[string](str)
	fmt.Printf("test: serializeContent() -> [buf:%v] [status:%v] [content:%v]\n", buf != nil, status, string(buf))

	str = "this is []byte content"
	buf, status = serializeContent[[]byte]([]byte(str))
	fmt.Printf("test: serializeContent() -> [buf:%v] [status:%v] [content:%v]\n", buf != nil, status, string(buf))

	str = "this is io.Reader content"
	r := strings.NewReader(str)
	buf, status = serializeContent[io.Reader](r)
	fmt.Printf("test: serializeContent() -> [buf:%v] [status:%v] [content:%v]\n", buf != nil, status, string(buf))

	str = "this is io.ReaderCloser content"
	r = strings.NewReader(str)
	buf, status = serializeContent[io.ReadCloser](io.NopCloser(r))
	fmt.Printf("test: serializeContent() -> [buf:%v] [status:%v] [content:%v]\n", buf != nil, status, string(buf))

	//Output:
	//test: serializeContent() -> [buf:false] [status:Internal [error: content type is invalid: int]]
	//test: serializeContent() -> [buf:true] [status:OK] [content:this is string content]
	//test: serializeContent() -> [buf:true] [status:OK] [content:this is []byte content]
	//test: serializeContent() -> [buf:true] [status:OK] [content:this is io.Reader content]
	//test: serializeContent() -> [buf:true] [status:OK] [content:this is io.ReaderCloser content]

}
