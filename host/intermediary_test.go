package host

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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
	ic := NewIntermediary(nil, nil)
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-nil-components -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewIntermediary(nil, serviceServeHTTP)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-service-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewIntermediary(authServeHTTP, serviceServeHTTP)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-auth-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ServeHTTP()-nil-components -> [status-code:200] [content:]
	//test: ServeHTTP()-service-only -> [status-code:200] [content:Service OK]
	//test: ServeHTTP()-auth-only -> [status-code:200] [content:Service OK]

}

func ExampleIntermediary_ServeHTTP() {
	ic := NewIntermediary(authServeHTTP, serviceServeHTTP)

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

func googleSearch(w http.ResponseWriter, r *http.Request) {
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	resp, _ := http.DefaultClient.Do(req)
	buf, err0 := io.ReadAll(resp.Body)
	if err0 != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	cnt, err := w.Write(buf)
	fmt.Printf("test: googleSearch() -> [cnt:%v] [err:%v] [error0:%v]\n", cnt, err, err0)
}

func ExampleNewControllerIntermediary() {
	im := NewControllerIntermediary("google-search", googleSearch)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: googleSearch() -> [cnt:110540] [err:<nil>] [status:OK]
	//test: NewControllerIntermediary() -> [status-code:200]

}

/*
func ExampleNewControllerIntermediary_5s() {
	im := NewControllerIntermediary("5s", "google-search", googleSearch)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Result().StatusCode)

	//Output:
	//test: googleSearch() -> [cnt:110540] [err:<nil>] [status:OK]
	//test: NewControllerIntermediary() -> [status-code:200]

}

func ExampleNewControllerIntermediary_100ms() {
	im := NewControllerIntermediary("100ms", "google-search", googleSearch)

	rec := exchange.NewResponseWriter() //httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=golang", nil)
	req.Header.Add(XRequestId, "1234-56-7890")
	req.Header.Add(XRelatesTo, "urn:business:activity")
	im(rec, req)
	fmt.Printf("test: NewControllerIntermediary() -> [status-code:%v]\n", rec.Response().StatusCode)

	//Output:
	//test: NewControllerIntermediary() -> [status-code:504]

}


*/