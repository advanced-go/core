package http2

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
)

type ContentWriter interface {
	Write(w http.ResponseWriter, content any)
}

type std struct{}

func (std) Write(w http.ResponseWriter, content any) {

}

// ContentWriterFn - function type for error handling
// type ErrorHandleFn func(requestId, location string, errs ...error) *Status
// NewErrorHandler - templated function providing an error handle function via a closure
type ContentWriterFunc func(http.ResponseWriter, any, runtime.Status, any)

/*
func NewContentWriter[E runtime.ErrorHandler](w http.ResponseWriter, content any, status runtime.Status, headers any) {
	//var e E
	//return func(w http.ResponseWriter, content any, status runtime.Status, headers any) {
	//return e.Handle(NewStatusError(http.StatusInternalServerError, location, errs...), requestId, "")
	WriteResponse[E, std](w, content, status, headers)
	//}
}


*/
/*
func NewContentWriter2(content any, r *http.Request) (w ContentWriter, contentEncoding string) {
	accept := r.Header.Get(contentEncoding)
	if len(accept) == 0 {
		var w std
		return
		return WriteResponse[E, std]
	}

}


*/
