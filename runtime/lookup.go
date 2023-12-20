package runtime

import (
	"fmt"
	"reflect"
)

//https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces

const (
	stringValueError = "error: stringFromType() value parameter is nil"
	listValueError   = "error: listFromType() value parameter is nil"
)

// LookupResultConstraints - lookup function constraints
type LookupResultConstraints interface {
	string | []string
}

type Lookup[T LookupResultConstraints] interface {
	Resolve(key string) (string, bool)
	SetOverride(t any)
}

type lookup[F LookupFunctionConstraints] struct {
	overrideFn F
	defaultFn  F
}

func (l *lookup[F]) SetOverride(t any) {
	if t == nil {
		l.overrideFn = nil
	}
}

func (l *lookup[F]) Resolve(key string) (string, bool) {
	return key, false
}

func NewLookup[T LookupResultConstraints, F LookupFunctionConstraints](defaultFn F) Lookup[T] {
	l := new(lookup[F])
	l.defaultFn = defaultFn
	return l
}

// LookupFunctionConstraints - lookup function constraints
type LookupFunctionConstraints interface {
	func(string) string | func(string) []string
}

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
