package host

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func serviceTestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Service OK")
}

func authTestHandler(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		tokenString := r.Header.Get(Authorization)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
		}
	}
}

func ExampleConditionalIntermediary_Nil() {
	ic := NewConditionalIntermediary(nil, nil, nil)
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ConditionalIntermediary()-nil-components -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(nil, serviceTestHandler, nil)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ConditionalIntermediary()-service-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(authTestHandler, serviceTestHandler, nil)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ConditionalIntermediary()-auth-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ConditionalIntermediary()-nil-components -> [status-code:500] [content:error: component 2 is nil]
	//test: ConditionalIntermediary()-service-only -> [status-code:200] [content:Service OK]
	//test: ConditionalIntermediary()-auth-only -> [status-code:200] [content:Service OK]

}

func ExampleConditionalIntermediary_HttpHandler() {
	ic := NewConditionalIntermediary(authTestHandler, serviceTestHandler, nil)

	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)

	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ConditionalIntermediary()-auth-failure -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")

	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ConditionalIntermediary()-auth-success -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ConditionalIntermediary()-auth-failure -> [status-code:401] [content:Missing authorization header]
	//test: ConditionalIntermediary()-auth-success -> [status-code:200] [content:Service OK]

}
