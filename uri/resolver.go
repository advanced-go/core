package uri

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

// ResolveFunc - type for resolution
type ResolveFunc func(string) string

// Resolver - resolver interface
type Resolver interface {
	SetOverride(t any, host string)
	Build(id string, values url.Values) string
}

// NewResolver - create a resolver
func NewResolver(defaultHost string, defaultFn ResolveFunc) Resolver {
	r := new(resolver)
	r.defaultHost = defaultHost
	r.defaultFn = defaultFn
	return r
}

type resolver struct {
	defaultHost string
	overrideFn  func(string) (string, string)
	defaultFn   ResolveFunc
}

// SetOverride - configure an override resolve func
func (r *resolver) SetOverride(t any, host string) {
	r.overrideFn = resolveFuncFromType(t, host)
}

// Build - perform resolution
func (r *resolver) Build(id string, values url.Values) string {
	override := false
	url := ""
	host := ""

	if r.overrideFn != nil {
		url, host = r.overrideFn(id)
		if len(url) > 0 {
			override = true
		}
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
		if len(host) == 0 {
			host = r.defaultHost
		}
		url = host + url
	}
	if !override && values != nil {
		url = url + "?" + values.Encode()
	}
	return url
}

func resolveFuncFromType(value any, host string) func(key string) (id string, host string) {
	if value == nil {
		if len(host) == 0 {
			return nil
		}
		return func(id string) (string, string) { return "", host }
	}
	switch ptr := value.(type) {
	case string:
		return func(id string) (string, string) { return ptr, host }
	case map[string]string:
		return func(id string) (string, string) {
			if v, ok := ptr[id]; ok {
				return v, host
			}
			return "", host
		}
	case func(string) (string, string):
		return ptr
	}
	return func(k string) (string, string) {
		return fmt.Sprintf("error: resolveFuncFromType() value parameter is an invalid type: %v", reflect.TypeOf(value)), host
	}
}
