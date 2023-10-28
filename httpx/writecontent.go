package httpx

import (
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"net/http"
	"reflect"
)

func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status *runtime.Status, location string) {
	var e E

	if status.Content() == nil {
		return
	}
	if w.Header().Get(ContentType) == "" {
		if status.JsonContent() {
			w.Header().Set(ContentType, ContentTypeJson)
		} else {
			w.Header().Set(ContentType, http.DetectContentType(status.Content()))
		}
	}
	w.Header().Set(ContentLength, fmt.Sprintf("%v", len(status.Content())))
	_, err := w.Write(status.Content())
	if err != nil {
		e.Handle(status.RequestId(), location+"/writeStatusContent", err)
	}
}

func serializeContent[T any](content T) ([]byte, *runtime.Status) {
	var buf []byte

	switch ptr := any(content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case io.Reader:
		var status *runtime.Status
		if ptr != nil {
			buf, status = ReadAll(io.NopCloser(ptr))
			if !status.OK() {
				return nil, status
			}
		}
	case io.ReadCloser:
		var status *runtime.Status
		if ptr != nil {
			buf, status = ReadAll(ptr)
			if !status.OK() {
				return nil, status
			}
		}
	default:
		return nil, runtime.NewStatusError(runtime.StatusInternal, serializeLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
	}
	return buf, runtime.NewStatusOK()
}
