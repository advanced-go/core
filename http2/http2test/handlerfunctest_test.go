package http2test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func ExampleHandlerFunc() {
	fn := HandlerFunc(healthHandler)
	fn.ServeHTTP(nil, nil)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "https://www.google.com/search?q=golang", nil)
	//genericHandler[healthHandler](rec, req)

	genericHandler[testServer](rec, req)

	fmt.Printf("test: HandlerFunc()")

	//Output:
}
