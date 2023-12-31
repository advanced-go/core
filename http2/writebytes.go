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
	bytesLoc = PkgPath + ":WriteBytes"
)

// WriteBytes -
func WriteBytes(content any, contentType string) ([]byte, string, runtime.Status) {
	var buf []byte

	switch ptr := (content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case error:
		buf = []byte(ptr.Error())
	case io.Reader:
		var status runtime.Status

		buf, status = runtime.NewBytes(ptr)
		if !status.OK() {
			return nil, "", status
		}
	case io.ReadCloser:
		var status runtime.Status

		buf, status = runtime.NewBytes(ptr)
		if !status.OK() {
			return nil, "", status
		}
	default:
		if strings.Contains(contentType, "json") {
			var err error

			buf, err = json.Marshal(content)
			if err != nil {
				status := runtime.NewStatusError(runtime.StatusJsonEncodeError, bytesLoc, err)
				if !status.OK() {
					return nil, "", status
				}
			}
			return buf, contentType, runtime.StatusOK()
		} else {
			return nil, "", runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
		}
	}
	if strings.Contains(contentType, "json") {
		return buf, ContentTypeJson, runtime.StatusOK()
	}
	return buf, http.DetectContentType(buf), runtime.StatusOK()
}
