package http2test

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}

func healthLivelinessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime.NewStatusOK()
	if status.OK() {
		//http2.WriteResponse[runtime.Output](w, []byte("up"), status, nil)
	} else {
		//http2.WriteResponse[runtime.Output](w, nil, status, nil)
	}
}
func example() {
	fn := HandlerFunc(healthLivelinessHandler)
	fn.ServeHTTP(nil, nil)
}
