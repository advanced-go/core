package host

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/http/httptest"
)

type authComponent struct {
}

func (ac *authComponent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r != nil {
		tokenString := r.Header.Get(Authorization)
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Missing authorization header")
		}
	}
}

type serviceComponent struct {
}

func (ac *serviceComponent) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Service OK")
}

func ExampleIntermediary_Nil() {
	auth := new(authComponent)
	serv := new(serviceComponent)

	ic := NewIntermediary(nil, nil)
	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic.ServeHTTP(rec, r)
	buf, _ := runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: ServeHTTP()-nil-components -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewIntermediary(nil, serv)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	ic.ServeHTTP(rec, r)
	buf, _ = runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: ServeHTTP()-service-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	ic = NewIntermediary(auth, serv)
	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")
	ic.ServeHTTP(rec, r)
	buf, _ = runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: ServeHTTP()-auth-only -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ServeHTTP()-nil-components -> [status-code:200] [content:]
	//test: ServeHTTP()-service-only -> [status-code:200] [content:Service OK]
	//test: ServeHTTP()-auth-only -> [status-code:200] [content:Service OK]
	
}

func ExampleIntermediary_ServeHTTP() {
	auth := new(authComponent)
	serv := new(serviceComponent)
	ic := NewIntermediary(auth, serv)

	rec := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)

	ic.ServeHTTP(rec, r)
	buf, _ := runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: ServeHTTP()-auth-failure -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	rec = httptest.NewRecorder()
	r, _ = http.NewRequest(http.MethodGet, "https://www.google.com/search?q-golang", nil)
	r.Header.Add(Authorization, "token")

	ic.ServeHTTP(rec, r)
	buf, _ = runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: ServeHTTP()-auth-success -> [status-code:%v] [content:%v]\n", rec.Result().StatusCode, string(buf))

	//Output:
	//test: ServeHTTP()-auth-failure -> [status-code:401] [content:Missing authorization header]
	//test: ServeHTTP()-auth-success -> [status-code:200] [content:Service OK]

}
