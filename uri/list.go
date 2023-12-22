package uri

import (
	"fmt"
	"reflect"
)

const (
	listValueError = "error: listFromType() value parameter is nil"
)

// ListLookupFunc - list lookup function
type ListLookupFunc func(string) []string

// ListLookup - lookup interface
type ListLookup interface {
	SetLookupFunc(value any)
	Lookup(key string) []string
}

type listLookup struct {
	lookupFn ListLookupFunc
}

// NewListLookup - new list lookup
func NewListLookup() ListLookup {
	return new(listLookup)
}

func (l *listLookup) SetLookupFunc(value any) {
	l.lookupFn = ListFromType(value)
}

func (l *listLookup) Lookup(key string) []string {
	if l.lookupFn == nil || len(key) == 0 {
		return nil
	}
	return l.lookupFn(key)
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
		return []string{fmt.Sprintf("error: listFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))}
	}
}
