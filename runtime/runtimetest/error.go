package runtimetest

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"github.com/go-ai-agent/core/strings"
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
	return h.HandleStatus(runtime.NewStatusError(runtime.StatusInternal, location, errs...), requestId, "")
}

func (h DebugError) HandleStatus(s *runtime.Status, requestId string, originUri string) *runtime.Status {
	if s == nil {
		return s
	}
	s.SetRequestId(requestId)
	s.SetOrigin(originUri)
	if s.IsErrors() && !s.ErrorsHandled() {
		fmt.Printf(defaultFormatter(s))
		s.SetErrorsHandled()
	}
	return s
}

func defaultFormatter(s *runtime.Status) string {
	return fmt.Sprintf("{ %v, %v, %v %v }\n",
		strings.JsonMarkup("id", s.RequestId(), true),
		strings.JsonMarkup("l", s.Location(), true),
		strings.JsonMarkup("o", s.Origin(), true),
		runtime.FormatErrors("err", s.Errors()))
}

/*
func formatErrors(errs []error) string {
	if len(errs) == 0 {
		return fmt.Sprintf("\"%v\" : null", "err")
	}
	result := fmt.Sprintf("\"%v\" : [ ", "err")
	for i, e := range errs {
		if i != 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%v\"", e.Error())
	}
	return result + " ]"
}


*/
