package runtimetest

import "fmt"
import "github.com/go-ai-agent/core/runtime"

// DebugError - debug error handler
type DebugError struct{}

func (h DebugError) Handle(requestId any, location string, errs ...error) *runtime.Status {
	if !runtime.IsErrors(errs) {
		return runtime.NewStatusOK()
	}
	return h.HandleStatus(runtime.NewStatusError(runtime.StatusInternal, location, errs...), requestId)
}

func (h DebugError) HandleStatus(s *runtime.Status, requestId any) *runtime.Status {
	if s != nil && s.IsErrors() && !s.ErrorsHandled() {
		loc := s.Location()
		if len(loc) == 0 {
			loc = runtime.EmptyArg
		}
		if s.RequestId() == "" {
			s.SetRequestId(requestId)
		}
		id := s.RequestId()
		if len(id) == 0 {
			id = runtime.EmptyArg
		}
		fmt.Printf("[%v %v %v]\n", id, loc, s.Errors())
		//s.RemoveErrors()
		s.ErrorsHandled()
	}
	return s
}
