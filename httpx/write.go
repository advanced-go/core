package httpx

import (
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

var writeLoc = PkgUri + "/write-response"

// WriteResponse - write a http.Response, utilizing the data, status, and headers for controlling the content
func WriteResponse[E runtime.ErrorHandler](w http.ResponseWriter, data []byte, status *runtime.Status, headersKV ...string) {
	writeResponse[E](w, data, status, headersKV...)
}

// WriteResponseCopy - write a http.Response, utilizing the data, status, and response for controlling the content
func WriteResponseCopy[E runtime.ErrorHandler](w http.ResponseWriter, resp *http.Response, headers ...string) {
	var buf []byte

	status := runtime.NewHttpStatusCode(resp.StatusCode)
	CreateHeaders(w.Header(), resp, headers...)
	if resp.Body != nil {
		var status1 *runtime.Status
		buf, status1 = ReadAll[E](resp.Body)
		if !status1.OK() {
			status = status1
		}
	}
	writeResponse[E](w, buf, status)
}

func writeResponse[E runtime.ErrorHandler](w http.ResponseWriter, data []byte, status *runtime.Status, headersKV ...string) {
	var e E
	if status == nil {
		status = runtime.NewStatusOK()
	}
	w.WriteHeader(status.Http())
	// if no data and there is content, then we need to set the ContentType
	if data == nil && status.Content() != nil {
		w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
	} else {
		SetHeaders(w, headersKV...)
	}
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
