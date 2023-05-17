package exchange

import (
	"github.com/go-sre/core/runtime"
	"net/http"
)

var writeLoc = PkgUrl + "/write-response"

// WriteResponse - write a http.Response, utilizing the data, status, and headers for controlling the content
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, data []byte, status *runtime.Status, headers ...string) {
	status.AddMetadata(w.Header(), headers...)
	writeResponse[E](w, data, status)
}

// WriteResponseCopy - write a http.Response, utilizing the data, status, and response for controlling the content
func WriteResponseCopy[E runtime.ErrorHandler](w http.ResponseWriter, resp *http.Response, status *runtime.Status, headers ...string) {
	var e E
	var buf []byte

	if status == nil {
		status = runtime.NewStatusOK()
	}
	CreateHeaders(w.Header(), resp, headers...)
	if status.OK() {
		var err error

		buf, err = ReadAll(resp.Body)
		if err != nil {
			status = e.Handle(nil, writeLoc, err)
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
	if status.OK() {
		if data != nil {
			_, ioErr = w.Write(data)
		}
	} else {
		if buf := status.Content(); buf != nil {
			_, ioErr = w.Write(buf)
		}
	}
	e.Handle(nil, writeLoc, ioErr)
}
