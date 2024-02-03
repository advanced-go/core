package http2test

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime.StatusOK()
	if status.OK() {
		//http2.WriteResponse[runtime.Output](w, []byte("up"), status, nil)
	} else {
		//http2.WriteResponse[runtime.Output](w, nil, status, nil)
	}
	fmt.Printf("test: in healthHandler.ServeHTTP()\n")
}

func genericHandler[T http.Handler](w http.ResponseWriter, r *http.Request) {
	var t T
	t.ServeHTTP(nil, nil)
	fmt.Printf("test: in genericHandler.ServeHTTP()\n")
}

type testServer struct{}

func (t testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("test: in testServer.ServeHTTP()\n")
}

type containerHandler struct {
	serveHttp func(w http.ResponseWriter, r *http.Request)
}

func NewContainerHandler(serve func(w http.ResponseWriter, r *http.Request)) *containerHandler {
	h := new(containerHandler)
	h.serveHttp = serve
	return h
}
