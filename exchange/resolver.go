package exchange

import (
	"fmt"
	"reflect"
	"strings"
)

// ResolveFunc - type for resolution
type ResolveFunc func(string) string

// Resolver - resolver interface
type Resolver interface {
	SetOverride(t any)
	Resolve(key string) string
}

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
func (r *resolver) Resolve(key string) string {
	uri := ""
	if r.overrideFn != nil {
		uri = r.overrideFn(key)
		if len(uri) > 0 {
			return uri
		}
	}
	uri = r.defaultFn(key)
	if len(uri) > 0 {
		return uri
	}
	if strings.HasPrefix(key, "/") {
		return r.defaultHost + "/" + uri
	}
	return key
}

func NewResolver(defaultHost string, defaultFn ResolveFunc) Resolver {
	r := new(resolver)
	r.defaultHost = defaultHost
	r.defaultFn = defaultFn
	return r
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
