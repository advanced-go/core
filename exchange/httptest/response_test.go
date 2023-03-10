package httptest

import (
	"fmt"
	"io"
)

func ExampleHtmlResponse() {
	resp, err0 := ReadResponse(fsys, "resource/http/html-response.html")
	if err0 != nil {
		fmt.Printf("test: ReadResponse() -> %v\n", err0)
	} else {
		fmt.Printf("test: ReadResponse() -> [resp:%v] [status:%v] [date:%v] [server:%v] [content-type:%v] [connection:%v]\n",
			resp != nil, resp.StatusCode, resp.Header.Get("Date"), resp.Header.Get("Server"), resp.Header.Get("Content-Type"), resp.Header.Get("Connection"))
		defer resp.Body.Close()
		buf, err := io.ReadAll(resp.Body)
		fmt.Printf("test: io.ReadAll() -> [error:%v] %v\n", err, string(buf))

	}

	//Output:
	//test: ReadResponse() -> [resp:true] [status:200] [date:Mon, 27 Jul 2009 12:28:53 GMT] [server:Apache/2.2.14 (Win32)] [content-type:text/html] [connection:Closed]
	//test: io.ReadAll() -> [error:<nil>] <html>
	//<body>
	//<h1>Hello, World!</h1>
	//</body>
	//</html>

}
