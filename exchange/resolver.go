package exchange

import (
	"fmt"
	"reflect"
	"strings"
)

type resolver struct {
	defaultHost string
	overrideFn  ResolveFunc
	defaultFn   ResolveFunc
}

// SetOverride - configure an override resolve func
func (r *resolver) SetOverride(t any) {
	r.overrideFn = resolveFuncFromType(t)
}

// Resolve - perform resolution
func (r *resolver) Resolve(id string) string {
	url := ""
	if r.overrideFn != nil {
		url = r.overrideFn(id)
	}
	if len(url) == 0 && r.defaultFn != nil {
		url = r.defaultFn(id)
	}
	if len(url) == 0 {
		url = id
	}
	if len(url) == 0 {
		return "error: id cannot be resolved to a URL"
	}
	if strings.HasPrefix(url, "/") {
		return r.defaultHost + url
	}
	return url
}

func resolveFuncFromType(value any) func(key string) string {
	if value == nil {
		return nil
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
		return fmt.Sprintf("error: resolveFuncFromType() value parameter is an invalid type: %v", reflect.TypeOf(value))
	}
}
