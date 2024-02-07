package uri

import (
	"fmt"
	"strings"
	"sync"
)

var (
	localAuthority = "localhost:8080"
)

const (
	httpScheme  = "http"
	httpsScheme = "https"
	localHost   = "localhost"
)

// Pair - key, value pair
type Pair struct {
	Key, Value string
}

// SetLocalAuthority - set the local authority
func SetLocalAuthority(authority string) {
	localAuthority = authority
}

// Resolver - type
type Resolver struct {
	override *sync.Map
}

// NewResolver - create a resolver
func NewResolver() *Resolver {
	return new(Resolver)
}

// SetOverrides - configure overrides
func (r *Resolver) SetOverrides(values []Pair) func() {
	if len(values) == 0 {
		r.override = nil
		return func() {}
	}
	r.override = new(sync.Map)
	for _, attr := range values {
		key, _ := TemplateToken(attr.Key)
		r.override.Store(key, attr.Value)
	}
	return func() {
		r.override = nil
	}
}

// Build - perform resolution
func (r *Resolver) Build(path string, values ...any) string {
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	return r.BuildWithAuthority(localAuthority, path, values...)
}

// BuildWithAuthority - perform resolution
func (r *Resolver) BuildWithAuthority(authority, path string, values ...any) string {
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	if r.override != nil {
		if uri, ok := r.OverrideUrl(path); ok {
			if len(values) > 0 && strings.Index(uri, "%v") != -1 {
				uri = fmt.Sprintf(uri, values...)
			}
			return uri
		}
	}
	if !strings.HasPrefix(path, "/") {
		path += "/"
	}
	if len(values) > 0 {
		path = fmt.Sprintf(path, values...)
	}
	scheme := httpsScheme
	if len(authority) == 0 || strings.HasPrefix(authority, localHost) {
		authority = localAuthority
		scheme = httpScheme
	}
	url2 := scheme + "://" + authority + path
	return url2
}

// OverrideUrl - return the overridden URL
func (r *Resolver) OverrideUrl(path string) (string, bool) {
	if r.override == nil {
		return "", false
	}
	if v, ok := r.override.Load(path); ok {
		if s, ok2 := v.(string); ok2 {
			return s, true
		}
	}
	return "", false
}
