package controller

import (
	"net/http"
	"sync/atomic"
)

const (
	PrimaryName   = "primary"
	SecondaryName = "secondary"
	primary       = 0
	secondary     = 1
)

type Resource struct {
	Name         string
	Authority    string `json:"authority"`
	LivenessPath string `json:"liveness"`
	handler      func(w *http.ResponseWriter, r *http.Request)
}

func (r *Resource) IsPrimary() bool {
	return r != nil && r.Name == PrimaryName
}

func (r *Resource) Uri(req *http.Request) string {
	return r.Authority + "/" + req.URL.Path
}

type Router struct {
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

func (r *Router) UpdateStats(resp *http.Response, rsc *Resource) {

}
