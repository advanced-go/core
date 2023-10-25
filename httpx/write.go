package httpx

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"net/http"
)

var writeLoc = PkgUri + "/write-response"

const (
	ContentLength = "Content-Length"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
func WriteResponse[E runtime.ErrorHandler, T any](w http.ResponseWriter, content T, status *runtime.Status, headersKV ...string) *runtime.Status {
	var e E
	var result error

	if status == nil {
		status = runtime.NewStatusOK()
	}
	// if missing a header value for a key, then write an internal error
	if (len(headersKV) & 1) == 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Handle(status.RequestId(), "WriteResponse", errors.New("invalid number of kv items: number is odd, possibly missing a value")).SetContent(http.StatusInternalServerError)
	}

	// always write status code and headers
	w.WriteHeader(status.Http())
	SetHeaders(w, headersKV...)

	// if status.Content is available, then that takes precedence
	if status.Content() != nil {
		if w.Header().Get(runtime.ContentType) == "" {
			w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
		}
		w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
		_, result = w.Write(status.Content())
	} else {
		switch ptr := any(content).(type) {
		case []byte:
			if ptr != nil {
				if w.Header().Get(ContentLength) == "" {
					w.Header().Set(ContentLength, fmt.Sprintf("%v", len(ptr)))
				}
				_, result = w.Write(ptr)
			}
		case string:
			buf := []byte(ptr)
			if buf != nil {
				if w.Header().Get(ContentLength) == "" {
					w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
				}
				_, result = w.Write(buf)
			}
		case io.ReadCloser:
			if ptr != nil {
				buf, status1 := ReadAll(ptr)
				status1.SetRequestId(status.RequestId())
				e.HandleStatus(status1)
				if status1.IsErrors() {
					result = status1.Errors()[0]
				} else {
					if w.Header().Get(ContentLength) == "" {
						w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
					}
					_, result = w.Write(buf)
				}
			}
		default:
			result = errors.New(fmt.Sprintf("error: content type is invalid [%v]", any(content)))
		}
	}
	return e.Handle(status.RequestId(), "WriteResponse", result)
}

// WriteMinResponse - write a http.Response, with status and headers and optional status content
func WriteMinResponse[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, headersKV ...string) *runtime.Status {
	var e E
	var result error

	if status == nil {
		status = runtime.NewStatusOK()
	}
	// if missing a header value for a key, then write an internal error
	if (len(headersKV) & 1) == 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Handle(status.RequestId(), "WriteResponse", errors.New("invalid number of kv items: number is odd, possibly missing a value")).SetContent(http.StatusInternalServerError)
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headersKV...)
	if status.Content() != nil {
		if w.Header().Get(runtime.ContentType) == "" {
			w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
		}
		w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
		_, result = w.Write(status.Content())
	}
	return e.Handle(status.RequestId(), "WriteNoContent", result)
}
