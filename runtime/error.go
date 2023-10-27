package runtime

import (
	"fmt"
	"github.com/go-ai-agent/core/strings"
	"log"
)

const (
	EmptyArg      = "[]"
	LocationName  = "Location"
	OriginName    = "Origin"
	RequestIdName = "RequestId"
	ErrorsName    = "Errors"
)

type FormatOutput func(s *Status) string

func SetFormatOutput(fn FormatOutput) {
	if fn != nil {
		formatter = fn
	}
}

var formatter FormatOutput = DefaultErrorFormatter

// ErrorHandleFn - function type for error handling
type ErrorHandleFn func(requestId any, location string, errs ...error) *Status

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(requestId any, location string, errs ...error) *Status
	HandleStatus(s *Status, requestId any, originUri string) *Status
}

// BypassError - bypass error handler
type BypassError struct{}

func (h BypassError) Handle(requestId any, _ string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatusError(StatusInternal, "", errs...)
}

func (h BypassError) HandleStatus(s *Status, _ any, _ string) *Status {
	return s
}

// LogError - debug error handler
type LogError struct{}

func (h LogError) Handle(requestId any, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatusError(StatusInternal, location, errs...), requestId, "")
}

func (h LogError) HandleStatus(s *Status, requestId any, originUri string) *Status {
	if s == nil {
		return s
	}
	if s.RequestId() == "" {
		s.SetRequestId(requestId)
	}
	if len(s.Origin()) == 0 && len(originUri) != 0 {
		s.SetOrigin(originUri)
	}
	if s != nil && s.IsErrors() && !s.ErrorsHandled() {
		log.Println(formatter(s))
		s.SetErrorsHandled()
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

func DefaultErrorFormatter(s *Status) string {
	return fmt.Sprintf("{ %v, %v, %v %v }\n",
		strings.JsonMarkup(RequestIdName, s.RequestId(), true),
		strings.JsonMarkup(LocationName, s.Location(), true),
		strings.JsonMarkup(OriginName, s.Origin(), true),
		FormatErrors(s.Errors()))
}

func FormatErrors(errs []error) string {
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
