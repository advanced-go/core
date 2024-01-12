package uri

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"strings"
	"sync"
)

const (
	SchemeHttps = "https"
	SchemeHttp  = "http"
	localHost   = "localhost"
)

var (
	localAuthority = "localhost:8080"
)

// SetLocalAuthority - set the local authority
func SetLocalAuthority(authority string) {
	localAuthority = authority
}

// Resolver - resolver interface
type Resolver interface {
	SetLocalHostOverride(v bool)
	SetAuthorities(values []runtime.Pair)
	SetOverrides(values []runtime.Pair)
	Build(authority, path string, values ...any) string
	Authority(authority string) (string, error)
	OverrideUrl(authority string) (string, bool)
}

// NewResolver - create a resolver
func NewResolver() Resolver {
	r := new(resolver)
	r.authority = new(sync.Map)
	return r
}

// NewResolverWithAuthorities - create a resolver with authorities
func NewResolverWithAuthorities(values []runtime.Pair) Resolver {
	r := new(resolver)
	r.authority = new(sync.Map)
	r.SetAuthorities(values)
	return r
}

type resolver struct {
	authority *sync.Map
	override  *sync.Map
	localHost bool
}

// SetLocalHostOverride - override authority to localhost
func (r *resolver) SetLocalHostOverride(v bool) {
	r.localHost = v
}

// SetAuthorities - configure authorities
func (r *resolver) SetAuthorities(values []runtime.Pair) {
	if len(values) == 0 {
		return
	}
	r.authority = new(sync.Map)
	for _, pair := range values {
		key, _ := TemplateToken(pair.Key)
		r.authority.Store(key, pair.Value)
	}
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
func (r *resolver) Build(authority, path string, values ...any) string {
	if len(authority) == 0 {
		return "resolver error: invalid argument, authority is empty"
	}
	if len(path) == 0 {
		return "resolver error: invalid argument, path is empty"
	}
	if r.override != nil {
		if u, ok := r.OverrideUrl(authority); ok {
			return u
		}
	}
	scheme := SchemeHttps
	if r.localHost {
		authority = localAuthority
		scheme = SchemeHttp
	} else {
		var err error
		authority, err = r.Authority(authority)
		if err != nil {
			return err.Error()
		}
		if strings.HasPrefix(authority, localHost) {
			scheme = SchemeHttp
		}
	}
	if len(values) > 0 {
		path = fmt.Sprintf(path, values...)
	}
	if !strings.HasPrefix(path, "/") {
		path += "/"
	}
	url2 := scheme + "://" + authority + path
	return url2
}

// Authority - return the authority using an expansion if needed
func (r *resolver) Authority(authority string) (string, error) {
	t, ok := TemplateToken(authority)
	if !ok {
		return authority, nil
	}
	if v, ok2 := r.authority.Load(t); ok2 {
		if s, ok3 := v.(string); ok3 {
			return s, nil
		}
	}
	return "", errors.New(fmt.Sprintf("resolver error: authority not found for variable: %v", authority))
}

// OverrideUrl - return the overridden URL
func (r *resolver) OverrideUrl(authority string) (string, bool) {
	t, ok := TemplateToken(authority)
	if !ok || r.override == nil {
		return "", false
	}
	if v, ok2 := r.override.Load(t); ok2 {
		if s, ok3 := v.(string); ok3 {
			return s, true
		}
	}
	return "", false
}
