package runtime

import (
	"fmt"
	"github.com/go-ai-agent/core/strings"
	"log"
)

const (
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
type ErrorHandleFn func(requestId, location string, errs ...error) *Status

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(requestId string, location string, errs ...error) *Status
	HandleStatus(s *Status, requestId string, originUri string) *Status
}

// BypassError - bypass error handler
type BypassError struct{}

func (h BypassError) Handle(requestId string, _ string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return NewStatusError(StatusInternal, "", errs...)
}

func (h BypassError) HandleStatus(s *Status, _ string, _ string) *Status {
	return s
}

// LogError - log error handler
type LogError struct{}

func (h LogError) Handle(requestId string, location string, errs ...error) *Status {
	if !IsErrors(errs) {
		return NewStatusOK()
	}
	return h.HandleStatus(NewStatusError(StatusInternal, location, errs...), requestId, "")
}

func (h LogError) HandleStatus(s *Status, requestId string, originUri string) *Status {
	if s == nil {
		return s
	}
	s.SetRequestId(requestId)
	s.SetOrigin(originUri)
	if s != nil && s.IsErrors() && !s.ErrorsHandled() {
		log.Println(formatter(s))
		s.SetErrorsHandled()
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

func DefaultErrorFormatter(s *Status) string {
	return fmt.Sprintf("{ %v, %v, %v %v }\n",
		strings.JsonMarkup(RequestIdName, s.RequestId(), true),
		strings.JsonMarkup(LocationName, s.Location(), true),
		strings.JsonMarkup(OriginName, s.Origin(), true),
		FormatErrors(ErrorsName, s.Errors()))
}

func FormatErrors(name string, errs []error) string {
	if len(errs) == 0 {
		return fmt.Sprintf("\"%v\" : null", name)
	}
	result := fmt.Sprintf("\"%v\" : [ ", name)
	for i, e := range errs {
		if i != 0 {
			result += ","
		}
		result += fmt.Sprintf("\"%v\"", e.Error())
	}
	return result + " ]"
}
