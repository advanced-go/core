package host

import "net/http"

const (
	Authorization = "Authorization"
)

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
	wrap := newWrapper(w)
	if i.c1 != nil {
		i.c1.ServeHTTP(wrap, r)
	}
	if wrap.statusCode == http.StatusOK && i.c2 != nil {
		i.c2.ServeHTTP(w, r)
	}
}

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
