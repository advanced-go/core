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
type Formatter func(code int, status string, errs []error, trace []string, attrs []any) string

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
type Logger func(code int, status string, errs []error, trace []string, attrs []any)

// SetErrorLogger - optional override of logging
func SetErrorLogger(fn Logger) {
	if fn != nil {
		logger = fn
	}
}

var (
	formatter            = defaultFormatter
	logger               = defaultLogger
	defaultLogger Logger = func(code int, status string, errs []error, trace []string, attrs []any) {
		log.Default().Println(formatter(code, status, errs, trace, attrs))
	}
)

// ErrorHandler - error handler interface
type ErrorHandler interface {
	Handle(s *Status, requestId ...string) *Status
}

// Bypass - bypass error handler
type Bypass struct{}

// Handle - bypass error handler
func (h Bypass) Handle(s *Status, _ ...string) *Status {
	return s
}

// Output - standard output error handler
type Output struct{}

// Handle - output error handler
func (h Output) Handle(s *Status, requestId ...string) *Status {
	if s == nil {
		return StatusOK()
	}
	if s.OK() {
		return s
	}
	if s.Error() != nil && !s.handled {
		s.addParentLocation()
		addRequestIdAttr(s, requestId)
		fmt.Printf("%v", formatter(s.Code, HttpStatus(s.Code), []error{s.Error()}, s.Trace(), s.Attrs()))
		s.handled = true
	}
	return s
}

// Log - log error handler
type Log struct{}

// Handle - log error handler
func (h Log) Handle(s *Status, requestId ...string) *Status {
	if s == nil {
		return StatusOK()
	}
	if s.OK() {
		return s
	}
	if s.Error() != nil && !s.handled {
		s.addParentLocation()
		addRequestIdAttr(s, requestId)
		go logger(s.Code, HttpStatus(s.Code), []error{s.Error()}, s.Trace(), s.Attrs())
		s.handled = true
	}
	return s
}

func defaultFormatter(code int, status string, errs []error, trace []string, attrs []any) string {
	str := strconv.Itoa(code)
	name := RequestIdName
	val := ""
	// TODO: add formatting of attributes
	if len(attrs) > 0 {
		if s, ok := attrs[0].(string); ok {
			val = s
		}
	}
	return fmt.Sprintf("{ %v, %v, %v, %v, %v }\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, status, true),
		jsonMarkup(name, val, true),
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
func OutputFormatter(code int, status string, errs []error, trace []string, attrs []any) string {
	str := strconv.Itoa(code)
	name := RequestIdName
	val := ""
	// TODO: add formatting of attributes
	if len(attrs) > 0 {
		if s, ok := attrs[0].(string); ok {
			val = s
		}
	}
	return fmt.Sprintf("{ %v, %v, %v, %v, %v\n}\n",
		jsonMarkup(StatusCodeName, str, false),
		jsonMarkup(StatusName, status, true),
		jsonMarkup(name, val, true),
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

// TODO: update when complete attribute handling is finished
func addRequestIdAttr(s *Status, requestId []string) {
	switch len(requestId) {
	case 1:
		s.AddAttr("", requestId[0])
	case 2:
		s.AddAttr("", requestId[1])
	}
}
