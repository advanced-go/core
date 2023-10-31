package runtimetest

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/strings"
	"net/http"
)

func SetFormatOutput(fn runtime.FormatOutput) {
	if fn != nil {
		formatter = fn
	}
}

var formatter runtime.FormatOutput = defaultFormatter

// DebugError - debug error handler
type DebugError struct{}

func (h DebugError) Handle(requestId string, location string, errs ...error) *runtime.Status {
	if !runtime.IsErrors(errs) {
		return runtime.NewStatusOK()
	}
	return h.HandleStatus(runtime.NewStatusError(http.StatusInternalServerError, location, errs...), requestId, "")
}

func (h DebugError) HandleStatus(s *runtime.Status, requestId string, location string) *runtime.Status {
	if s == nil {
		return s
	}
	s.SetRequestId(requestId)
	s.AddLocation(location)
	if s.IsErrors() && !s.ErrorsHandled() {
		fmt.Printf(defaultFormatter(s))
		s.SetErrorsHandled()
	}
	return s
}

func defaultFormatter(s *runtime.Status) string {
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		strings.JsonMarkup("c", fmt.Sprintf("%v", s.Code()), false),
		strings.JsonMarkup("s", s.Description(), true),
		strings.JsonMarkup("id", s.RequestId(), true),
		runtime.FormatLocation("l", s.Location()),
		runtime.FormatErrors("err", s.Errors()))
}
