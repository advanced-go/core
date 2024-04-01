package host

import (
	"context"
	"fmt"
	"github.com/advanced-go/core/access"
	"github.com/advanced-go/core/controller"
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

func ExampleNewIngressControllerIntermediary_Nil() {
	access.EnableInternalLogging()
	im := NewIngressControllerIntermediary(nil, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewIngressControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressControllerIntermediary() -> [status-code:200]

}

func ExampleNewIngressControllerIntermediary_5s() {
	ctrl := controller.NewTimeoutController("google-search", time.Second*5)
	im := NewIngressControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewIngressControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressControllerIntermediary() -> [status-code:200]

}

func ExampleNewIngressControllerIntermediary_1ms() {
	ctrl := controller.NewTimeoutController("google-search", time.Millisecond*1)
	im := NewIngressControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewIngressControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>] [write-err:<nil>]
	//test: NewIngressControllerIntermediary() -> [status-code:504]

}

func ExampleNewIngressControllerIntermediary_100ms() {
	ctrl := controller.NewTimeoutController("google-search", time.Millisecond*100)
	im := NewIngressControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewIngressControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:<nil>] [read-err:context deadline exceeded] [write-err:<nil>]
	//test: NewIngressControllerIntermediary() -> [status-code:504]

}
