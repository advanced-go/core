package runtime

import (
	"fmt"
	"reflect"
)

const (
	//stringValueError = "error: stringFromType() value parameter is nil"
	listValueError = "error: ListFromType() value parameter is nil"
)

/*
// LookupFunctionConstraints - lookup function constraints
type LookupFunctionConstraints interface {
	func(string) string | func(string) []string
}

// LookupFromType - create a function from a type
func LookupFromType[F LookupFunctionConstraints](t any) (r F) {
	switch ptr := any(&r).(type) {
	case *func(string) string:
		*ptr = StringFromType(t)
	case *func(string) []string:
		*ptr = ListFromType(t)
	}
	return r
}

func StringFromType(value any) func(key string) string {
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

*/

// ListLookupFunc - list lookup function
type ListLookupFunc func(string) []string

// NewListLookup - new list lookup from a value type
func NewListLookup(value any) ListLookupFunc {
	if value == nil {
		return nil
	}
	return ListFromType(value)
}

// ListLookup - lookup function
func ListLookup(key string, fn ListLookupFunc) []string {
	if fn == nil || len(key) == 0 {
		return nil
	}
	val := fn(key)
	if len(val) > 0 {
		return val
	}
	return nil
}

// ListFromType - create a function returning the value
func ListFromType(value any) ListLookupFunc {
	if value == nil {
		return func(key string) []string { return []string{listValueError} }
	}
	if s, ok := value.(string); ok {
		return func(key string) []string { return []string{s, ""} }
	}
	if l, ok := value.([]string); ok {
		return func(key string) []string { return l }
	}
	if m, ok := value.(map[string][]string); ok {
		return func(key string) []string {
			if v, ok1 := m[key]; ok1 {
				return v
			}
			return nil
		}
	}
	if fn, ok := value.(func(string) []string); ok {
		return fn
	}
	return func(key string) []string {
		return []string{fmt.Sprintf("error: ListFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))}
	}
}
