package host

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleIntermediary_Nil() {
	ic := NewConditionalIntermediary(nil, nil, nil)
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ := io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-nil-components -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(nil, serviceTestHandler, nil)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic(rec, r)
	buf, _ = io.ReadAll(rec.Result().Body)
	fmt.Printf("test: ServeHTTP()-service-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewConditionalIntermediary(authTestHandler, serviceTestHandler, nil)
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
