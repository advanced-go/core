package uri

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
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

// SetLocalAuthority - set the local authority
func SetLocalAuthority(authority string) {
	localAuthority = authority
}

// Resolver - resolver interface
type Resolver interface {
	SetOverrides(values []runtime.Pair)
	Build(path string, values ...any) string
	BuildWithAuthority(authority, path string, values ...any) string
	OverrideUrl(key string) (string, bool)
}

// NewResolver - create a resolver
func NewResolver() Resolver {
	r := new(resolver)
	return r
}

type resolver struct {
	override *sync.Map
}

// SetOverrides - configure overrides
func (r *resolver) SetOverrides(values []runtime.Pair) {
	if len(values) == 0 {
		r.override = nil
		return
	}
	r.override = new(sync.Map)
	for _, attr := range values {
		key, _ := TemplateToken(attr.Key)
		r.override.Store(key, attr.Value)
	}
}

// Build - perform resolution
func (r *resolver) Build(path string, values ...any) string {
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	return r.BuildWithAuthority(localAuthority, path, values...)
	/*
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
		url2 := "http://" + localAuthority + path
		return url2

	*/
}

// BuildWithAuthority - perform resolution
func (r *resolver) BuildWithAuthority(authority, path string, values ...any) string {
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
	url2 := scheme + "//" + authority + path
	return url2
}

// OverrideUrl - return the overridden URL
func (r *resolver) OverrideUrl(path string) (string, bool) {
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
