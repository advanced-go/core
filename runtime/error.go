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

// FormatOutput - output formatting type
type FormatOutput func(s Status) string

// SetFormatOutput - optional override of output formatting
func SetFormatOutput(fn FormatOutput) {
	if fn != nil {
		formatter = fn
	}
}

var formatter FormatOutput = defaultErrorFormatter

// ErrorHandler - template parameter error handler interface
type ErrorHandler interface {
	Handle(s Status, requestId string, callerLocation string) Status
}

// Bypass - bypass error handler
type Bypass struct{}

func (h Bypass) Handle(s Status, _ string, _ string) Status {
	return s
}

// Output - standard output error handler
type Output struct{}

func (h Output) Handle(s Status, requestId string, location string) Status {
	if s == nil {
		return StatusOK()
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

// Log - log error handler
type Log struct{}

func (h Log) Handle(s Status, requestId string, callerLocation string) Status {
	if s == nil {
		return StatusOK()
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

// ErrorHandleFn - function type for error handling
//type ErrorHandleFn func(requestId, location string, errs ...error) *Status
// NewErrorHandler - templated function providing an error handle function via a closure
/*
func NewErrorHandler[E ErrorHandler]() ErrorHandleFn {
	var e E
	return func(requestId string, location string, errs ...error) *Status {
		return e.Handle(NewStatusError(http.StatusInternalServerError, location, errs...), requestId, "")
	}
}
*/

func defaultErrorFormatter(s Status) string {
	str := strconv.Itoa(s.Code())
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		strings.JsonMarkup(StatusCodeName, str, false),
		strings.JsonMarkup(StatusName, s.Description(), true),
		strings.JsonMarkup(RequestIdName, s.RequestId(), true),
		formatTrace(TraceName, s.Location()),
		formatErrors(ErrorsName, s.Errors()))
}

func formatErrors(name string, errs []error) string {
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

func formatTrace(name string, trace []string) string {
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

// NewInvalidBodyTypeError - invalid type error
func NewInvalidBodyTypeError(t any) error {
	return errors.New(fmt.Sprintf("invalid body type: %v", reflect.TypeOf(t)))
}
