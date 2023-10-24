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
		return e.Handle(nil, "WriteResponse", errors.New("invalid number of kv items: number is odd, possibly missing a value")).SetContent(http.StatusInternalServerError)
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
			if w.Header().Get(ContentLength) == "" {
				w.Header().Set(ContentLength, fmt.Sprintf("%v", len(ptr)))
			}
			_, result = w.Write([]byte(ptr))
		case io.ReadCloser:
			if ptr != nil {
				buf, status1 := ReadAll[E](ptr)
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
	return e.Handle(nil, "WriteResponse", result)
}

// WriteResponseNoContent - write a http.Response, utilizing status and headers
func WriteResponseNoContent[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status) *runtime.Status {
	var e E
	var result error

	if status == nil {
		status = runtime.NewStatusOK()
	}
	w.WriteHeader(status.Http())
	if status.Content() != nil {
		if w.Header().Get(runtime.ContentType) == "" {
			w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
		}
		w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
		_, result = w.Write(status.Content())
	}
	return e.Handle(nil, "WriteNoContent", result)
}
