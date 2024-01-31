package controller

import (
	"fmt"
	"net/http"
)

func ExampleCreateRequest() {
	req, _ := http.NewRequest(http.MethodPut, "https://www.google.com/search?q=golang", nil)

	fmt.Printf("test: CreateRequest() -> [method:%v] [uri:%v]\n", req.Method, req.URL.String())

	//Output:
}
