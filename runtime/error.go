package runtime

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/strings"
	"log"
	"reflect"
	"strconv"
)

const (
	StatusName     = "status"
	StatusCodeName = "code"
	TraceName      = "trace"
	RequestIdName  = "request-id"
	ErrorsName     = "errors"
)

type FormatOutput func(s Status) string

func SetFormatOutput(fn FormatOutput) {
	if fn != nil {
		formatter = fn
	}
}

var formatter FormatOutput = DefaultErrorFormatter

// ErrorHandleFn - function type for error handling
//type ErrorHandleFn func(requestId, location string, errs ...error) *Status

// ErrorHandler - template parameter error handler interface
// Handle(requestId string, location string, errs ...error) *Status
type ErrorHandler interface {
	Handle(s Status, requestId string, callerLocation string) Status
}

// BypassError - bypass error handler
type BypassError struct{}

func (h BypassError) Handle(s Status, _ string, _ string) Status {
	return s
}

// DebugError - debug error handler
type DebugError struct{}

func (h DebugError) Handle(s Status, requestId string, location string) Status {
	if s == nil || s.OK() {
		return NewStatusOK()
	}
	if s.OK() {
		return s
	}
	s.SetRequestId(requestId)
	s.AddLocation(location)
	if s.IsErrors() {
		log.Println(formatter(s))
		setErrorsHandled(s)
	}
	return s
}

// LogError - log error handler
type LogError struct{}

func (h LogError) Handle(s Status, requestId string, callerLocation string) Status {
	if s == nil {
		return NewStatusOK()
	}
	if s.OK() {
		return s
	}
	s.SetRequestId(requestId)
	s.AddLocation(callerLocation)
	if s.IsErrors() && !errorsHandled(s) {
		log.Println(formatter(s))
		setErrorsHandled(s)
	}
	return s
}

// NewErrorHandler - templated function providing an error handle function via a closure
/*
func NewErrorHandler[E ErrorHandler]() ErrorHandleFn {
	var e E
	return func(requestId string, location string, errs ...error) *Status {
		return e.Handle(NewStatusError(http.StatusInternalServerError, location, errs...), requestId, "")
	}
}


*/

func DefaultErrorFormatter(s Status) string {
	str := strconv.Itoa(s.Code())
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		strings.JsonMarkup(StatusCodeName, str, false),
		strings.JsonMarkup(StatusName, s.Description(), true),
		strings.JsonMarkup(RequestIdName, s.RequestId(), true),
		FormatTrace(TraceName, s.Location()),
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

func FormatTrace(name string, trace []string) string {
	if len(trace) == 0 {
		return fmt.Sprintf("\"%v\" : null", name)
	}
	result := fmt.Sprintf("\"%v\" : [ ", name)
	for i := len(trace) - 1; i >= 0; i-- {
		if i < len(trace)-1 {
			result += ","
		}
		result += fmt.Sprintf("\"%v\"", trace[i])
	}
	return result + " ]"
}

func NewInvalidBodyTypeError(t any) error {
	return errors.New(fmt.Sprintf("invalid body type: %v", reflect.TypeOf(t)))
}
