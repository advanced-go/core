package host

import "net/http"

type Intermediary interface {
}

type intermediary struct {
	c1, c2 http.Handler
}

func NewIntermediary(c1 http.Handler, c2 http.Handler) http.Handler {
	i := new(intermediary)
	i.c1 = c1
	i.c2 = c2
	return i
}

func (i *intermediary) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	i.c1.ServeHTTP(w, r)
	//if w.resp == nil ||
}
