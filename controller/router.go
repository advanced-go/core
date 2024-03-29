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

type Router struct {
	activeHost         atomic.Int64
	HealthLivenessPath string `json:"liveness"`
	PrimaryHost        string `json:"primary"`
	SecondaryHost      string `json:"secondary"`
}

func (r *Router) Uri(path string) string {
	if r.activeHost.Load() == primary {
		return r.PrimaryHost + path
	} else {
		return r.SecondaryHost + path
	}
}

func (r *Router) swapHost() (swapped bool) {
	old := r.activeHost.Load()
	if old == primary {
		swapped = r.activeHost.CompareAndSwap(old, secondary)
	} else {
		swapped = r.activeHost.CompareAndSwap(old, primary)
	}
	return
}

func (r *Router) PingHost() string {
	if r.activeHost.Load() == primary {
		return r.SecondaryHost
	} else {
		return r.PrimaryHost
	}
}

func (r *Router) Liveness() (statusCode int) {
	//r,_ := http.NewRequest(http.MethodGet,r.PingHost() + r.HealthLivenessPath,nil)
	return http.StatusOK
}
