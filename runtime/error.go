package runtime

import (
	"log"
)

const (
	EmptyArg = "[]"
)

// ErrorHandleFn - function type for error handling
type ErrorHandleFn func(requestId any, location string, errs ...error) *Status

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(requestId any, location string, errs ...error) *Status
	HandleStatus(s *Status, requestId any) *Status
}

// BypassError - bypass error handler
type BypassError struct{}

func (h BypassError) Handle(requestId any, _ string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatusError(StatusInternal, "", errs...)
}

func (h BypassError) HandleStatus(s *Status, _ any) *Status {
	return s
}

// LogError - debug error handler
type LogError struct{}

func (h LogError) Handle(requestId any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatusError(StatusInternal, location, errs...), requestId)
}

func (h LogError) HandleStatus(s *Status, requestId any) *Status {
	if s != nil && s.IsErrors() && !s.ErrorsHandled() {
		loc := ifElse(s.Location(), EmptyArg)
		if s.RequestId() == "" {
			s.SetRequestId(requestId)
		}
		req := ifElse(s.RequestId(), EmptyArg)
		log.Println(req, loc, s.Errors())
		//s.RemoveErrors()
		s.SetErrorsHandled()
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
	return func(requestId any, location string, errs ...error) *Status {
		return e.Handle(requestId, location, errs...)
	}
}
