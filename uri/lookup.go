package uri

import (
	"fmt"
	"reflect"
)

const (
	stringValueError = "error: stringFromType() value parameter is nil"
)

// LookupFunc - lookup function
type LookupFunc func(string) string

// Lookup - lookup interface
type Lookup interface {
	SetFunc(value any)
	Value(key string) (string, bool)
}

type uriLookup struct {
	fn LookupFunc
}

// NewLookup - new uri lookup
func NewLookup() Lookup {
	return new(uriLookup)
}

// SetFunc - set the lookup function
func (l *uriLookup) SetFunc(value any) {
	l.fn = stringFromType(value)
}

// Value - return the value associated with the key
func (l *uriLookup) Value(key string) (string, bool) {
	if l.fn == nil || len(key) == 0 {
		return "", false
	}
	val := l.fn(key)
	if len(val) > 0 {
		return val, true
	}
	return "", false
}

func stringFromType(value any) func(key string) string {
	if value == nil {
		return func(k string) string { return stringValueError }
	}
	switch ptr := value.(type) {
	case string:
		return func(k string) string { return ptr }
	case map[string]string:
		return func(k string) string {
			v := ptr[k]
			if len(v) > 0 {
				return v
			}
			return ""
		}
	case func(string) string:
		return ptr
	}
	return func(k string) string {
		return fmt.Sprintf("error: stringFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))
	}
}
