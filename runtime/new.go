package runtime

import (
	"encoding/json"
	"errors"
	"fmt"
	uri2 "github.com/advanced-go/core/uri"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

const (
	StatusOKUri       = "urn:status:ok"
	StatusNotFoundUri = "urn:status:notfound"
	StatusTimeoutUri  = "urn:status:timeout"
	newLoc            = PkgPath + ":New"
)

// New - create a new type from any JSON content
func New[T any](v any) (t T, status Status) {
	var buf []byte

	switch ptr := v.(type) {
	case string:
		if uri2.IsStatusURL(ptr) {
			return t, NewStatusFrom(ptr)
		}
		buf, status = NewBytes(ptr)
		if !status.OK() {
			return
		}
		if ptr1, ok := any(&t).(*[]byte); ok {
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return
	case *url.URL:
		if uri2.IsStatusURL(ptr.String()) {
			return t, NewStatusFrom(ptr.String())
		}
		buf, status = NewBytes(ptr.String())
		if !status.OK() {
			return
		}
		if ptr1, ok := any(&t).(*[]byte); ok {
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return
	case []byte:
		buf = ptr
		if ptr1, ok := any(&t).(*[]byte); ok {
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return
	case *http.Response:
		if ptr1, ok := any(&t).(*[]byte); ok {
			buf, status = NewBytes(ptr)
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
	case io.Reader:
		if ptr1, ok := any(&t).(*[]byte); ok {
			buf, status = NewBytes(ptr)
			if !status.OK() {
				return
			}
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.NewDecoder(ptr).Decode(&t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()
	case io.ReadCloser:
		if ptr1, ok := any(&t).(*[]byte); ok {
			buf, status = NewBytes(ptr)
			if !status.OK() {
				return
			}
			*ptr1 = buf
			return t, StatusOK()
		}
		err := json.NewDecoder(ptr).Decode(&t)
		_ = ptr.Close()
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()
	default:
		return t, NewStatusError(StatusInvalidArgument, newLoc, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
	}
}
