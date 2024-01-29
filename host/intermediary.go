package host

import "net/http"

const (
	Authorization = "Authorization"
)

//type intermediary struct {
//	c1, c2 http.HandlerFunc
//}

type ServeHTTPFunc func(w http.ResponseWriter, r *http.Request)

func NewIntermediary(c1 ServeHTTPFunc, c2 ServeHTTPFunc) ServeHTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrap := newWrapper(w)
		if c1 != nil {
			c1(wrap, r)
		}
		if wrap.statusCode == http.StatusOK && c2 != nil {
			c2(w, r)
		}
	}
}

/*
func (i *intermediary) serveHTTP(w http.ResponseWriter, r *http.Request) {
	wrap := newWrapper(w)
	if i.c1 != nil {
		i.c1.ServeHTTP(wrap, r)
	}
	if wrap.statusCode == http.StatusOK && i.c2 != nil {
		i.c2.ServeHTTP(w, r)
	}
}


*/

type wrapper struct {
	writer     http.ResponseWriter
	statusCode int
}

func newWrapper(writer http.ResponseWriter) *wrapper {
	w := new(wrapper)
	w.writer = writer
	w.statusCode = http.StatusOK
	return w
}

func (w *wrapper) Header() http.Header {
	return w.writer.Header()
}

func (w *wrapper) Write(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *wrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
}
