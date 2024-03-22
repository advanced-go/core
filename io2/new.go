package io2

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

const (
	StatusOKUri       = "urn:status:ok"
	StatusNotFoundUri = "urn:status:notfound"
	StatusTimeoutUri  = "urn:status:timeout"
	newLoc            = runtime.PkgPath + ":New"
)

// New - create a new type from JSON content, supporting: string, *url.URL, []byte, io.Reader, io.ReadCloser
// Note: content encoded []byte is not supported
func New[T any](v any, h http.Header) (t T, status *runtime.Status) {
	var buf []byte

	switch ptr := v.(type) {
	case string:
		if isStatusURL(ptr) {
			return t, NewStatusFrom(ptr)
		}
		buf, status = ReadFile(ptr)
		if !status.OK() {
			return
		}
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
		}
		return
	case *url.URL:
		if isStatusURL(ptr.String()) {
			return t, NewStatusFrom(ptr.String())
		}
		buf, status = ReadFile(ptr.String())
		if !status.OK() {
			return
		}
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
		}
		return
	case []byte:
		// TO DO : determine if encoding is supported for []byte
		buf = ptr
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
		}
		return
	case io.Reader:
		reader, status0 := NewEncodingReader(ptr, h)
		if !status0.OK() {
			return t, status0.AddLocation()
		}
		err := json.NewDecoder(reader).Decode(&t)
		_ = reader.Close()
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
		}
		return t, runtime.StatusOK()
	case io.ReadCloser:
		reader, status0 := NewEncodingReader(ptr, h)
		if !status0.OK() {
			return t, status0.AddLocation()
		}
		err := json.NewDecoder(reader).Decode(&t)
		_ = reader.Close()
		_ = ptr.Close()
		if err != nil {
			return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, err)
		}
		return t, runtime.StatusOK()
	default:
		return t, runtime.NewStatusError(runtime.StatusInvalidArgument, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
	}
}

/*
	case *http.Response:
		if ptr1, ok := any(&t).(*[]byte); ok {
			buf, status = ReadAll(ptr.Body,h)
			if !status.OK() {
				return
			}
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.NewDecoder(ptr.Body).Decode(&t)
		_ = ptr.Body.Close()
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()
	case *http.Request:
		if ptr1, ok := any(&t).(*[]byte); ok {
			buf, status = ReadAll(ptr.Body)
			if !status.OK() {
				return
			}
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.NewDecoder(ptr.Body).Decode(&t)
		_ = ptr.Body.Close()
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()

*/
