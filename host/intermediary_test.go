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

func authServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		tokenString := r.Header.Get(Authorization)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
		}
	}
}

func serviceServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Service OK")
}

func ExampleIntermediary_Nil() {
	ic := NewConditionalIntermediary(nil, nil, nil)
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-nil-components -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(nil, serviceServeHTTP, nil)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-service-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(authServeHTTP, serviceServeHTTP, nil)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-auth-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ServeHTTP()-nil-components -> [status-code:500] [content:error: component 2 is nil]
	//test: ServeHTTP()-service-only -> [status-code:200] [content:Service OK]
	//test: ServeHTTP()-auth-only -> [status-code:200] [content:Service OK]

}

func ExampleIntermediary_ServeHTTP() {
	ic := NewConditionalIntermediary(authServeHTTP, serviceServeHTTP, nil)

	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)

	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-auth-failure -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")

	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-auth-success -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ServeHTTP()-auth-failure -> [status-code:401] [content:Missing authorization header]
	//test: ServeHTTP()-auth-success -> [status-code:200] [content:Service OK]

}

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

func ExampleNewControllerIntermediary_Nil() {
	access.EnableInternalLogging()
	im := NewControllerIntermediary(nil, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewControllerIntermediary() -> [status-code:200]

}

func ExampleNewControllerIntermediary_5s() {
	ctrl := new(controller.Controller)
	ctrl.RouteName = "google-search"
	ctrl.Timeout.Duration = time.Second * 5
	im := NewControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:true] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: NewControllerIntermediary() -> [status-code:200]

}

func ExampleNewControllerIntermediary_1ms() {
	ctrl := new(controller.Controller)
	ctrl.RouteName = "google-search"
	ctrl.Timeout.Duration = time.Millisecond * 1
	im := NewControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>] [write-err:<nil>]
	//test: NewControllerIntermediary() -> [status-code:504]

}

func ExampleNewControllerIntermediary_100ms() {
	ctrl := new(controller.Controller)
	ctrl.RouteName = "google-search"
	ctrl.Timeout.Duration = time.Millisecond * 900
	im := NewControllerIntermediary(ctrl, httpCall)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add("X-Request-Id", "1234-56-7890")
	req.Header.Add("X-Relates-To", "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: httpCall() -> [content:false] [do-err:<nil>] [read-err:context deadline exceeded] [write-err:<nil>]
	//test: NewControllerIntermediary() -> [status-code:504]

}
