package http2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"reflect"
	"strings"
)

const (
	AcceptEncoding = "Accept-Encoding"
	writeLoc       = PkgPath + ":WriteResponse"
	gzipEncoding   = "gzip"
	bytesLoc       = PkgPath + ":Bytes"
	jsonToken      = "json"
	contentType    = "Content-Type"
)

// WriteResponse - write a http.Response, utilizing the content, status, and headers
// Content types supported: []byte, string, error, io.Reader, io.ReadCloser. Other types will be treated as JSON and serialized, if
// the headers content type is JSON. If not JSON, then an error will be raised.
func WriteResponse[E runtime.ErrorHandler, W ContentWriter](w http.ResponseWriter, content any, status runtime.Status, headers any) {
	var e E

	if status == nil {
		status = runtime.StatusOK()
	}
	SetHeaders(w, headers)
	accept := w.Header().Get(AcceptEncoding)
	w.Header().Del(AcceptEncoding)
	w.WriteHeader(status.Http())
	status0 := writeContent(w, content, accept)
	if !status0.OK() {
		e.Handle(status0, status0.RequestId(), writeLoc)
	}
	/*
		ct := GetContentType(headers)
		buf, status0 := runtime.Bytes(content, ct)
		if !status0.OK() {
			e.Handle(status0, status.RequestId(), writeLoc)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		SetHeaders(w, headers)
		if len(ct) == 0 {
			w.Header().Set(ContentType, http.DetectContentType(buf))
		}
		w.WriteHeader(status.Http())

	*/
	/*
		bytes, err := w.Write(buf)
		if err != nil {
			e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, err), status.RequestId(), "")
		}
		if bytes != len(buf) {
			e.Handle(runtime.NewStatusError(http.StatusInternalServerError, writeLoc, errors.New(fmt.Sprintf("error on ResponseWriter().Write() -> [got:%v] [want:%v]\n", bytes, len(buf)))), status.RequestId(), "")
		}

	*/
	return
}

func writeContent(w io.Writer, content any, accept string) runtime.Status {
	var buf []byte
	gzip := false
	if strings.Contains(accept, gzipEncoding) {
		gzip = true
	}
	switch ptr := (content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case error:
		buf = []byte(ptr.Error())
	case io.Reader:
		var status runtime.Status

		buf, status = runtime.ReadAll(ptr, nil)
		if !status.OK() {
			return status
		}
	case io.ReadCloser:
		var status runtime.Status

		buf, status = runtime.ReadAll(ptr, nil)
		_ = ptr.Close()
		if !status.OK() {
			return status
		}
	default:
		if strings.Contains(contentType, jsonToken) {
			var err error

			buf, err = json.Marshal(content)
			if err != nil {
				status := runtime.NewStatusError(runtime.StatusJsonEncodeError, bytesLoc, err)
				if !status.OK() {
					return status
				}
			}
			return runtime.StatusOK()
		} else {
			return runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
		}
	}
	if buf != nil {
	}
	if gzip {
	}
	return runtime.StatusOK()
}

/*
func writeStatusContent[E runtime.ErrorHandler](w http.ResponseWriter, status runtime.Status, location string) {
	var e E

	if status.Content() == nil {
		return
	}
	ct := status.ContentHeader().Get(ContentType)
	buf, status1 := runtime.Bytes(status.Content(), ct)
	if !status1.OK() {
		e.Handle(status, status.RequestId(), location+writeStatusContentLoc)
		return
	}
	if len(ct) == 0 {
		ct = http.DetectContentType(buf)
	}
	w.Header().Set(ContentType, ct)
	_, err := w.Write(buf)
	if err != nil {
		e.Handle(runtime.NewStatusError(http.StatusInternalServerError, location+writeStatusContentLoc, err), "", "")
	}
}


*/
