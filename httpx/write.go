package httpx

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"net/http"
)

var (
	writeLoc    = PkgUri + "/write-response"
	minWriteLoc = PkgUri + "/write-min-response"
)

const (
	ContentLength = "Content-Length"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
func WriteResponse[E runtime.ErrorHandler, T any](w http.ResponseWriter, content T, status *runtime.Status, headersKV ...string) *runtime.Status {
	var e E

	if status == nil {
		status = runtime.NewStatusOK()
	}
	err := ValidateKVHeaders(headersKV...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Handle(status, writeLoc, err)
	}

	// always write status code and headers
	w.WriteHeader(status.Http())
	SetHeaders(w, headersKV...)

	// if status.Content is available, then that takes precedence
	if status.Content() != nil {
		return writeStatusContent[E](w, status, writeLoc)
	}
	switch ptr := any(content).(type) {
	case []byte:
		if ptr == nil {
			return runtime.NewStatusOK()
		}
		if w.Header().Get(ContentLength) == "" {
			w.Header().Set(ContentLength, fmt.Sprintf("%v", len(ptr)))
		}
		_, result := w.Write(ptr)
		return e.Handle(status, writeLoc, result)
	case string:
		if ptr == "" {
			return runtime.NewStatusOK()
		}
		buf := []byte(ptr)
		if w.Header().Get(ContentLength) == "" {
			w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
		}
		_, result := w.Write(buf)
		return e.Handle(status, writeLoc, result)
	case io.ReadCloser:
		if ptr == nil {
			return runtime.NewStatusOK()
		}
		buf, status1 := ReadAll(ptr)
		if !status1.OK() {
			return e.HandleStatus(status1.SetRequestId(status))
		}
		if w.Header().Get(ContentLength) == "" {
			w.Header().Set(ContentLength, fmt.Sprintf("%v", len(buf)))
		}
		_, result := w.Write(buf)
		return e.Handle(status, writeLoc, result)
	default:
		return e.Handle(status, writeLoc, errors.New(fmt.Sprintf("error: content type is invalid [%v]", any(content))))
	}
	//return runtime.NewStatusOK()
}

// WriteMinResponse - write a http.Response, with status, optional headers and optional status content
func WriteMinResponse[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, headersKV ...string) *runtime.Status {
	var e E

	if status == nil {
		status = runtime.NewStatusOK()
	}
	err := ValidateKVHeaders(headersKV...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return e.Handle(status, writeLoc, err)
	}
	w.WriteHeader(status.Http())
	SetHeaders(w, headersKV...)
	return writeStatusContent[E](w, status, minWriteLoc)
}

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, location string) *runtime.Status {
	var e E

	if status.Content() == nil {
		return runtime.NewStatusOK()
	}
	if w.Header().Get(runtime.ContentType) == "" {
		w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
	}
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
	_, result := w.Write(status.Content())
	return e.Handle(status, location, result)
}

//if w.Header().Get(runtime.ContentType) == "" {
//	w.Header().Set(runtime.ContentType, http.DetectContentType(status.Content()))
//}
//w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
//_, result := w.Write(status.Content())
