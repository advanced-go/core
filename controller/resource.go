package controller

import (
	"net/http"
)

type Resource struct {
	internal     bool
	Name         string `json:"name"`
	Authority    string `json:"authority"`
	LivenessPath string `json:"liveness"`
	handler      func(w http.ResponseWriter, r *http.Request)
	do           func(*http.Request) (*http.Response, error)
}

func NewResource(name, authority, path string, handler any) *Resource {
	r := new(Resource)
	r.internal = false
	r.Name = name
	r.Authority = authority
	r.LivenessPath = path
	if h, ok := handler.(func(http.ResponseWriter, *http.Request)); ok {
		r.internal = true
		r.handler = h
	}
	// http.DefaultClient.Do
	if h, ok := handler.(func(*http.Request) (*http.Response, error)); ok {
		r.do = h
	}
	// Leave handler nil and default
	return r
}

func (r *Resource) IsPrimary() bool {
	return r != nil && r.Name == PrimaryName
}

func (r *Resource) Uri(req *http.Request) string {
	return r.Authority + "/" + req.URL.Path
}
