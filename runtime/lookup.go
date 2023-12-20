package runtime

import (
	"fmt"
	"reflect"
)

const (
	stringValueError = "error: stringFromType() value parameter is nil"
	listValueError   = "error: listFromType() value parameter is nil"
)

// LookupFunctionConstraints - lookup function constraints
type LookupFunctionConstraints interface {
	func(string) string | func(string) []string
}

// LookupFromType - create a function from a type
func LookupFromType[F LookupFunctionConstraints](t any) (r F) {
	switch ptr := any(&r).(type) {
	case *func(string) string:
		*ptr = stringFromType(t)
	case *func(string) []string:
		*ptr = listFromType(t)
	}
	return r
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
			return k
		}
	case func(string) string:
		return ptr
	}
	return func(k string) string {
		return fmt.Sprintf("error: stringFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))
	}
}

func listFromType(value any) func(key string) []string {
	if value == nil {
		return func(key string) []string { return []string{listValueError} }
	}
	if s, ok := value.(string); ok {
		return func(key string) []string { return []string{s, ""} }
	}
	if s, ok := value.([]string); ok {
		return func(key string) []string { return s }
	}
	if m, ok := value.(map[string][]string); ok {
		return func(key string) []string { return m[key] }
	}
	if fn, ok := value.(func(string) []string); ok {
		return fn
	}
	return func(key string) []string {
		return []string{fmt.Sprintf("error: listFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))}
	}
}
