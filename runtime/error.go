package runtime

import (
	"fmt"
	"log"
)

const (
	emptyArg = "[]"
)

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(ctx any, location string, errs ...error) *Status
	//HandleStatus(ctx any, s *Status) *Status
}

// DebugError - debug error handler
type DebugError struct{}

func (h DebugError) Handle(ctx any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(ctx, NewStatus(StatusInternal, location, errs...).SetRequestId(ContextRequestId(ctx)))
}

func (h DebugError) HandleStatus(_ any, s *Status) *Status {
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

func (h LogError) Handle(ctx any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(ctx, NewStatus(StatusInternal, location, errs...).SetRequestId(ContextRequestId(ctx)))
}

func (h LogError) HandleStatus(_ any, s *Status) *Status {
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
