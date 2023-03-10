package exchange

import (
	"github.com/gotemplates/core/runtime"
	"net/http"
)

// WriteResponse - write a http.Response, utilizing the data, status, and headers for controlling the content
func WriteResponse(w http.ResponseWriter, data []byte, status *runtime.Status, headers ...string) {
	if status == nil {
		w.WriteHeader(http.StatusOK)
		if data != nil {
			w.Write(data)
		}
		return
	}
	status.AddMetadata(w.Header(), headers...)
	w.WriteHeader(status.Http())
	if status.OK() {
		if data != nil {
			w.Write(data)
		}
	} else {
		if buf := status.Content(); buf != nil && status != nil {
			w.Write(buf)
		}
	}
}
