package controller

import (
	"github.com/advanced-go/core/runtime"
	"net/http"
	"sync/atomic"
)

const (
	PrimaryName   = "primary"
	SecondaryName = "secondary"
	primary       = 0
	secondary     = 1
)

type Router struct {
	Do        func(r *http.Request) (*http.Response, *runtime.Status)
	primary   *Resource
	secondary *Resource
	active    atomic.Int64
}

func NewRouter(primary, secondary *Resource) *Router {
	r := new(Router)
	r.primary = primary
	r.primary.Name = PrimaryName
	r.secondary = secondary
	r.secondary.Name = SecondaryName
	return r
}

func (r *Router) RouteTo() *Resource {
	if r.active.Load() == primary {
		return r.primary
	}
	return r.secondary
}

func (r *Router) UpdateStats(statusCode int, rsc *Resource) {

}

func (r *Router) Swap() (swapped bool) {
	old := r.active.Load()
	if old == primary {
		swapped = r.active.CompareAndSwap(old, secondary)
	} else {
		swapped = r.active.CompareAndSwap(old, primary)
	}
	return
}
