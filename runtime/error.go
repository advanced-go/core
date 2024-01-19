package runtime

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	StatusName     = "status"
	StatusCodeName = "code"
	TraceName      = "trace"
	RequestIdName  = "request-id"
	ErrorsName     = "errors"
	githubHost     = "github"
	githubDotCom   = "github.com"
	githubTemplate = "https://%v/tree/main%v"
	fragmentId     = "#"
)

// Formatter - output formatting type
type Formatter func(s Status) string

// SetErrorFormatter - optional override of output formatting
func SetErrorFormatter(fn Formatter) {
	if fn != nil {
		formatter = fn
	}
}

// SetOutputFormatter - optional override of output formatting
func SetOutputFormatter() {
	SetErrorFormatter(OutputFormatter)
}

// Logger - log function
type Logger func(s Status)

// SetErrorLogger - optional override of logging
func SetErrorLogger(fn Logger) {
	if fn != nil {
		logger = fn
	}
}

var (
	formatter            = defaultFormatter
	logger               = defaultLogger
	defaultLogger Logger = func(s Status) { log.Default().Println(formatter(s)) }
)

// ErrorHandler - error handler interface
type ErrorHandler interface {
	Handle(s Status, requestId string, callerLocation string) Status
}

// Bypass - bypass error handler
type Bypass struct{}

// Handle - bypass error handler
func (h Bypass) Handle(s Status, _ string, _ string) Status {
	return s
}

// Output - standard output error handler
type Output struct{}

// Handle - output error handler
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
		fmt.Printf("%v", formatter(s))
		setErrorsHandled(s)
	}
	return s
}

// Log - log error handler
type Log struct{}

// Handle - log error handler
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
		logger(s) //log.Println(formatter(s))
		setErrorsHandled(s)
	}
	return s
}

func defaultFormatter(s Status) string {
	str := strconv.Itoa(s.Code())
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, s.Description(), true),
		jsonMarkup(RequestIdName, s.RequestId(), true),
		formatTrace(TraceName, s.Location()),
		formatErrors(ErrorsName, s.ErrorList()))
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
		result += fmt.Sprintf("\"%v\"", formatUri(trace[i]))
	}
	return result + " ]"
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

// OutputFormatter - formatter for special output formatting
func OutputFormatter(s Status) string {
	str := strconv.Itoa(s.Code())
	return fmt.Sprintf("{ %v, %v, %v, %v, %v\n}\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, s.Description(), true),
		jsonMarkup(RequestIdName, s.RequestId(), true),
		outputFormatTrace(TraceName, s.Location()),
		outputFormatErrors(ErrorsName, s.ErrorList()))
}

func outputFormatErrors(name string, errs []error) string {
	if len(errs) == 0 {
		return fmt.Sprintf("\"%v\" : null", name)
	}
	result := fmt.Sprintf("\n\"%v\" : [\n", name)
	for i, e := range errs {
		if i != 0 {
			result += ",\n"
		}
		result += fmt.Sprintf("  \"%v\"", e.Error())
	}
	return result + " \n]"
}

func outputFormatTrace(name string, trace []string) string {
	if len(trace) == 0 {
		return fmt.Sprintf("\"%v\" : null", name)
	}
	result := fmt.Sprintf("\n\"%v\" : [\n", name)
	for i := len(trace) - 1; i >= 0; i-- {
		if i < len(trace)-1 {
			result += ",\n"
		}
		result += fmt.Sprintf("  \"%v\"", formatUri(trace[i]))
	}
	return result + " \n]"
}

func formatUri(uri string) string {
	i := strings.Index(uri, githubHost)
	if i == -1 {
		return uri
	}
	uri = strings.Replace(uri, githubHost, githubDotCom, len(githubDotCom))
	i = strings.LastIndex(uri, "/")
	if i != -1 {
		first := uri[:i]
		last := uri[i:]
		last = strings.Replace(last, urnSeparator, fragmentId, len(fragmentId))
		return fmt.Sprintf(githubTemplate, first, last)
	}
	return uri
}

// NewInvalidBodyTypeError - invalid type error
func NewInvalidBodyTypeError(t any) error {
	return errors.New(fmt.Sprintf("invalid body type: %v", reflect.TypeOf(t)))
}

const (
	markupNull   = "\"%v\":null"
	markupString = "\"%v\":\"%v\""
	markupValue  = "\"%v\":%v"
)

// jsonMarkup - markup a name/value pair
func jsonMarkup(name, value string, stringValue bool) string {
	if len(value) == 0 {
		return fmt.Sprintf(markupNull, name)
	}
	format := markupString
	if !stringValue {
		format = markupValue
	}
	return fmt.Sprintf(format, name, value)
}
