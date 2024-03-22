package exchange

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/url"
)

func readAll(body io.ReadCloser) ([]byte, *runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, err)
	}
	return buf, runtime.StatusOK()
}

func Example_ReadResponse() {
	s := "file://[cwd]/exchangetest/test-response.txt"
	u, _ := url.Parse(s)

	resp, status0 := readResponse(u)
	fmt.Printf("test: readResponse(%v) -> [status:%v] [statusCode:%v]\n", s, status0, resp.StatusCode)

	buf, status := readAll(resp.Body)
	fmt.Printf("test: readAll() -> [status:%v] [content-length:%v]\n", status, len(buf)) //string(buf))

	//Output:
	//test: readResponse(file://[cwd]/exchangetest/test-response.txt) -> [status:OK] [statusCode:200]
	//test: readAll() -> [status:OK] [content-length:58]

}

func Example_ReadResponse_URL_Nil() {
	resp, status0 := readResponse(nil)
	fmt.Printf("test: readResponse(nil) -> [error:[%v]] [statusCode:%v]\n", status0.Error(), resp.StatusCode)

	//Output:
	//test: readResponse(nil) -> [error:[error: URL is nil]] [statusCode:500]

}

func Example_ReadResponse_Invalid_Scheme() {
	s := "https://www.google.com/search?q=golang"
	u, _ := url.Parse(s)

	resp, status0 := readResponse(u)
	fmt.Printf("test: readResponse(%vl) -> [error:[%v]] [statusCode:%v]\n", s, status0.Error(), resp.StatusCode)

	//Output:
	//test: readResponse(https://www.google.com/search?q=golangl) -> [error:[error: Invalid URL scheme : https]] [statusCode:500]

}

func Example_ReadResponse_HTTP_Error() {
	s := "file://[cwd]/exchangetest/message.txt"
	u, _ := url.Parse(s)

	resp, status0 := readResponse(u)
	fmt.Printf("test: readResponse(%v) -> [error:[%v]] [statusCode:%v]\n", s, status0.Error(), resp.StatusCode)

	//Output:
	//test: readResponse(file://[cwd]/exchangetest/message.txt) -> [error:[malformed HTTP status code "text"]] [statusCode:500]

}

func Example_ReadResponse_NotFound() {
	s := "file://[cwd]/exchangetest/not-found.txt"
	u, _ := url.Parse(s)

	resp, status0 := readResponse(u)
	fmt.Printf("test: readResponse(%v) -> [error:[%v]] [statusCode:%v]\n", s, status0.Error(), resp.StatusCode)

	//Output:
	//test: readResponse(file://[cwd]/exchangetest/not-found.txt) -> [error:[open C:\Users\markb\GitHub\core\exchange\exchangetest\not-found.txt: The system cannot find the file specified.]] [statusCode:404]

}

func Example_ReadResponse_EOF_Error() {
	s := "file://[cwd]/exchangetest/http-503-error.txt"
	u, _ := url.Parse(s)

	resp, status0 := readResponse(u)
	fmt.Printf("test: readResponse(%v) -> [error:[%v]] [statusCode:%v]\n", s, status0.Error(), resp.StatusCode)

	//Output:
	//test: readResponse(file://[cwd]/exchangetest/http-503-error.txt) -> [error:[unexpected EOF]] [statusCode:500]

}
