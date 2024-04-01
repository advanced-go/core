package host

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"io"
	"net/http"
	"net/http/httptest"
	"time"
)

func httpCall(w http.ResponseWriter, r *http.Request) {
	cnt := 0
	var err2 error
	var err1 error
	var buf []byte

	resp, err0 := http.DefaultClient.Do(r)
	if err0 != nil {
		if r.Context().Err() == context.DeadlineExceeded {
			w.WriteHeader(http.StatusGatewayTimeout)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		buf, err1 = io.ReadAll(resp.Body)
		if err1 != nil {
			if err1 == context.DeadlineExceeded {
				w.WriteHeader(http.StatusGatewayTimeout)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			cnt, err2 = w.Write(buf)
			w.WriteHeader(http.StatusOK)
		}
	}
	fmt.Printf("test: httpCall() -> [content:%v] [do-err:%v] [read-err:%v] [write-err:%v]\n", cnt > 0, err0, err1, err2)
}

func ExampleNewIngressTimeoutIntermediary_Nil() {
	access.EnableInternalLogging()
	im := NewIngressTimeoutIntermediary("google-search", 0, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewIngressTimeoutIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressTimeoutIntermediary() -> [status-code:200]

}

func ExampleNewIngressTimeoutIntermediary_5s() {
	im := NewIngressTimeoutIntermediary("google-search", time.Second*5, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewIngressTimeoutIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressTimeoutIntermediary() -> [status-code:200]

}

func ExampleNewIngressTimeoutIntermediary_1ms() {
	im := NewIngressTimeoutIntermediary("google-search", time.Millisecond*1, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewIngressTimeoutIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressTimeoutIntermediary() -> [status-code:504]

}

func ExampleNewIngressTimeoutIntermediary_100ms() {
	im := NewIngressTimeoutIntermediary("google-search", time.Millisecond*100, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewIngressTimeoutIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:<nil>] [read-err:context deadline exceeded] [write-err:<nil>]
	//test: NewIngressTimeoutIntermediary() -> [status-code:504]

}
