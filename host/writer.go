package host

import (
	"net/http"
	"time"
)

type wrapper struct {
	writer     http.ResponseWriter
	statusCode int
	written    int64
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
	w.written += int64(len(p))
	return w.writer.Write(p)
}

func (w *wrapper) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.writer.WriteHeader(statusCode)
}

// Milliseconds - convert time.Duration to milliseconds
func Milliseconds(duration time.Duration) int {
	return int(duration / time.Duration(1e6))
}
