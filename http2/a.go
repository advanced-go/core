package http2

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

func init() {
	fmt.Println("a.go -> init()")
}

func example() {

}

type Exchange func(r *http.Request) (*http.Response, runtime.Status)

var do Exchange

func DoA(r *http.Request) (*http.Response, runtime.Status) {
	return do(r)
}
