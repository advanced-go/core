package runtime

import (
	"fmt"
	"log"
)

const (
	emptyArg = "[]"
)

// ErrorHandler2 - template parameter error handler interface
type ErrorHandler2 interface {
	Handle(ctx any, location string, errs ...error) *Status
	//HandleStatus(ctx any, s *Status) *Status
}

// DebugError2 - debug error handler
type DebugError2 struct{}

func (h DebugError2) Handle(ctx any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(ctx, NewStatus(StatusInternal, location, errs...).SetRequestId(ContextRequestId(ctx)))
}

func (h DebugError2) HandleStatus(_ any, s *Status) *Status {
	if s != nil && s.IsErrors() {
		loc := ifElse(s.Location(), emptyArg)
		req := ifElse(s.RequestId(), emptyArg)
		fmt.Printf("[%v %v %v]\n", req, loc, s.Errors())
		s.RemoveErrors()
	}
	return s
}

// LogError2 - debug error handler
type LogError2 struct{}

func (h LogError2) Handle(ctx any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(ctx, NewStatus(StatusInternal, location, errs...).SetRequestId(ContextRequestId(ctx)))
}

func (h LogError2) HandleStatus(_ any, s *Status) *Status {
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
