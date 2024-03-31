package controller

import (
	"github.com/advanced-go/core/runtime"
	"golang.org/x/time/rate"
	"net/http"
)

type Resource2 struct {
	doLimiter   *rate.Limiter
	pingLimiter *rate.Limiter
}

func NewResource() *Resource2 {
	r := new(Resource2)
	r.doLimiter = rate.NewLimiter(100, 10)
	r.pingLimiter = rate.NewLimiter(100, 10)

	return r
}

func (r *Resource2) Ping() *runtime.Status {
	return runtime.StatusOK()
}

func (r *Resource2) Do(req *http.Request) (*http.Response, *runtime.Status) {
	return nil, runtime.StatusOK()
}
