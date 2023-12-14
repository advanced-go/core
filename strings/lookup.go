package strings

import (
	"fmt"
	"reflect"
)

var (
	overrideLookup func(string) string
	defaultLookup  = errorLookup
	errorLookup    = func(key string) string { return "error: no default Lookup has been initialized" }
)

func setDefaultLookup(t any) {
	newDefault := LookupFromType(t)
	if newDefault == nil {
		defaultLookup = func(key string) string {
			return fmt.Sprintf("error: invalid default Lookup type: %v", reflect.TypeOf(t))
		}
	} else {
		defaultLookup = newDefault
	}
}

func setOverrideLookup(t any) {
	if t == nil {
		overrideLookup = nil
		return
	}
	overrideLookup = LookupFromType(t)
	if overrideLookup == nil {
		overrideLookup = func(key string) string {
			return fmt.Sprintf("error: invalid override Lookup type: %v", reflect.TypeOf(t))
		}
	}
}

func lookup(key string) string {
	if overrideLookup != nil {
		val := overrideLookup(key)
		if len(val) > 0 {
			return val
		}
	}
	return defaultLookup(key)
}

func LookupFromType(t any) func(string) string {
	switch ptr := t.(type) {
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
	return nil
}
