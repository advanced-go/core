package exchange

import (
	"github.com/go-ai-agent/core/httpx"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

var writeLoc = PkgUri + "/write-response"

// WriteResponse - write a http.Response, utilizing the data, status, and headers for controlling the content
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, data []byte, status *runtime.Status, headers ...string) {
	status.AddMetadata(w.Header(), headers...)
	writeResponse[E](w, data, status)
}

// WriteResponseCopy - write a http.Response, utilizing the data, status, and response for controlling the content
func WriteResponseCopy[E runtime.ErrorHandler](w http.ResponseWriter, resp *http.Response, headers ...string) {
	var buf []byte

	status := runtime.NewHttpStatusCode(resp.StatusCode)
	httpx.CreateHeaders(w.Header(), resp, headers...)
	if resp.Body != nil {
		var status1 *runtime.Status
		buf, status1 = httpx.ReadAll[E](resp.Body)
		if !status1.OK() {
			status = status1
		}
	}
	writeResponse[E](w, buf, status)
}

func writeResponse[E runtime.ErrorHandler](w http.ResponseWriter, data []byte, status *runtime.Status) {
	var e E
	if status == nil {
		status = runtime.NewStatusOK()
	}
	w.WriteHeader(status.Http())
	var ioErr error
	if data != nil {
		_, ioErr = w.Write(data)
	} else {
		if buf := status.Content(); buf != nil {
			_, ioErr = w.Write(buf)
		}
	}
	e.Handle(nil, writeLoc, ioErr)
}
