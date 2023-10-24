package runtime

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const (
	emptyArg = "[]"
)

// ErrorHandleFn - function type for error handling
type ErrorHandleFn func(requestId string, location string, errs ...error) *Status

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(requestId string, location string, errs ...error) *Status
}

// BypassError - bypass error handler
type BypassError struct{}

func (h BypassError) Handle(requestId string, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatus(StatusInternal, errs...) //.SetRequestId(requestId)
}

//func (h BypassError) HandleStatus(_ any, s *Status) *Status {
//	return s
//}

// DebugError - debug error handler
type DebugError struct{}

func (h DebugError) Handle(requestId, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, errs...).SetRequestId(requestId).SetLocation(location))
}

func (h DebugError) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		loc := ifElse(s.Location(), emptyArg)
		req := ifElse(s.RequestId(), emptyArg)
		fmt.Printf("[%v %v %v]\n", req, loc, s.Errors())
		s.RemoveErrors()
	}
	return s
}

// LogError - debug error handler
type LogError struct{}

func (h LogError) Handle(requestId, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatus(StatusInternal, errs...).SetRequestId(requestId).SetLocation(location))
}

func (h LogError) HandleStatus(s *Status) *Status {
	if s != nil && s.IsErrors() {
		loc := ifElse(s.Location(), emptyArg)
		req := ifElse(s.RequestId(), emptyArg)
		log.Println(req, loc, s.Errors())
		s.RemoveErrors()
	}
	return s
}

func ifElse(s string, def string) string {
	if len(s) == 0 {
		return def
	}
	return s
}

// NewErrorHandler - templated function providing an error handle function via a closure
func NewErrorHandler[E ErrorHandler]() ErrorHandleFn {
	var e E
	return func(requestId string, location string, errs ...error) *Status {
		return e.Handle(requestId, location, errs...)
	}
}

func RequestId(t any) string {
	if ctx, ok := t.(context.Context); ok {
		return ContextRequestId(ctx)
	}
	if req, ok := t.(*http.Request); ok {
		return req.Header.Get(XRequestIdName)
	}
	return ""
}
