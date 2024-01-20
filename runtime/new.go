package runtime

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
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

// New - create a new type from JSON content, supporting: string, *url.URL, []byte, io.Reader, io.ReadCloser
// Note: content encoded []byte is not supported
func New[T any](v any, h http.Header) (t T, status Status) {
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
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
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
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return
	case []byte:
		buf = ptr
		err := json.Unmarshal(buf, &t)
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return
	case io.Reader:
		var err error

		encoding := contentEncoding(h)
		switch encoding {
		case GzipEncoding:
			zr, err1 := gzip.NewReader(ptr)
			if err1 != nil {
				return t, NewStatusError(StatusGzipDecodingError, readAllLoc, err1)
			}
			err = json.NewDecoder(zr).Decode(&t)
			_ = zr.Close()
		case NoneEncoding:
			err = json.NewDecoder(ptr).Decode(&t)
		default:
			return t, NewStatusError(StatusContentEncodingError, readAllLoc, encodingError(encoding))
		}
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()
	case io.ReadCloser:
		var err error
		encoding := contentEncoding(h)
		switch encoding {
		case GzipEncoding:
			zr, err1 := gzip.NewReader(ptr)
			if err1 != nil {
				return t, NewStatusError(StatusGzipDecodingError, readAllLoc, err1)
			}
			err = json.NewDecoder(zr).Decode(&t)
			_ = zr.Close()
		case NoneEncoding:
			err = json.NewDecoder(ptr).Decode(&t)
		default:
			return t, NewStatusError(StatusContentEncodingError, readAllLoc, encodingError(encoding))
		}
		_ = ptr.Close()
		if err != nil {
			return t, NewStatusError(StatusJsonDecodeError, newLoc, err)
		}
		return t, StatusOK()
	default:
		return t, NewStatusError(StatusInvalidArgument, newLoc, errors.New(fmt.Sprintf("error: invalid type [%v]", reflect.TypeOf(v))))
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
