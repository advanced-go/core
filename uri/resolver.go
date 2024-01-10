package uri

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// TO DO : Need to configure a way to set default authority
const (
	SchemeHttps = "https"
	SchemeHttp  = "http"
	localHost   = "localhost"
)

var (
	localAuthority = "localhost:8080"
)

type Attr struct {
	Key, Value string
}

func SetLocalAuthority(authority string) {
	localAuthority = authority
}

// Resolver - resolver interface
type Resolver interface {
	SetLocalHostOverride(v bool)
	SetAuthorities(values []Attr)
	SetOverrides(values []Attr)
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
func NewResolverWithAuthorities(values []Attr) Resolver {
	r := new(resolver)
	r.authority = new(sync.Map)
	r.SetAuthorities(values)
	//for _, attr := range authorities {
	//	r.authority.Store(attr.Key, attr.Value)
	//}
	return r
}

type resolver struct {
	authority *sync.Map
	override  *sync.Map
	localHost bool
}

func (r *resolver) SetLocalHostOverride(v bool) {
	r.localHost = v
}

// SetAuthorities - configure authorities
func (r *resolver) SetAuthorities(values []Attr) {
	if len(values) == 0 {
		return
	}
	r.authority = new(sync.Map)
	for _, attr := range values {
		r.authority.Store(attr.Key, attr.Value)
	}
}

// SetOverrides - configure overrides
func (r *resolver) SetOverrides(values []Attr) {
	if len(values) == 0 {
		r.override = nil
		return
	}
	r.override = new(sync.Map)
	for _, attr := range values {
		r.override.Store(attr.Key, attr.Value)
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
	url2 := scheme + "://" + authority + path
	return url2
}

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
