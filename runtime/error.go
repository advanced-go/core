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
	urnSeparator   = ":"
)

// Formatter - output formatting type
type Formatter func(code int, errs []error, trace []string, content any, requestId string) string

// SetErrorFormatter - optional override of error formatting
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
type Logger func(code int, errs []error, trace []string, content any, requestId string)

// SetErrorLogger - optional override of logging
func SetErrorLogger(fn Logger) {
	if fn != nil {
		logger = fn
	}
}

var (
	formatter            = defaultFormatter
	logger               = defaultLogger
	defaultLogger Logger = func(code int, errs []error, trace []string, content any, requestId string) {
		log.Default().Println(formatter(code, errs, trace, content, requestId))
	}
)

// ErrorHandler - error handler interface
type ErrorHandler interface {
	Handle(s *Status, requestId string) *Status
}

// Bypass - bypass error handler
type Bypass struct{}

// Handle - bypass error handler
func (h Bypass) Handle(s *Status, _ string) *Status {
	return s
}

// Output - standard output error handler
type Output struct{}

// Handle - output error handler
func (h Output) Handle(s *Status, requestId string) *Status {
	if s == nil {
		return StatusOK()
	}
	if s.OK() {
		return s
	}
	if s.Error() != nil && !s.handled {
		s.addLocation(3)
		fmt.Printf("%v", formatter(s.Code, []error{s.Error()}, s.Trace(), s.Content(), requestId))
		s.handled = true
	}
	return s
}

// Log - log error handler
type Log struct{}

// Handle - log error handler
func (h Log) Handle(s *Status, requestId string) *Status {
	if s == nil {
		return StatusOK()
	}
	if s.OK() {
		return s
	}
	if s.Error() != nil && !s.handled {
		s.addLocation(3)
		logger(s.Code, []error{s.Error()}, s.Trace(), s.Content(), requestId)
		s.handled = true
	}
	return s
}

func defaultFormatter(code int, errs []error, trace []string, _ any, requestId string) string {
	str := strconv.Itoa(code)
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, HttpStatus(code), true),
		jsonMarkup(RequestIdName, requestId, true),
		formatErrors(ErrorsName, errs),
		formatTrace(TraceName, trace))
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
	if len(errs) == 0 || errs[0] == nil {
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
func OutputFormatter(code int, errs []error, trace []string, _ any, requestId string) string {
	str := strconv.Itoa(code)
	return fmt.Sprintf("{ %v, %v, %v, %v, %v\n}\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, HttpStatus(code), true),
		jsonMarkup(RequestIdName, requestId, true),
		outputFormatErrors(ErrorsName, errs),
		outputFormatTrace(TraceName, trace))

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
