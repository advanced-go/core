package http2

import (
	"errors"
	"fmt"
	io2 "github.com/go-ai-agent/core/io"
	"github.com/go-ai-agent/core/json"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var (
	bytesLoc = PkgUri + "/WriteBytes"
)

func WriteBytes(content any, contentType string) ([]byte, string, *runtime.Status) {
	var buf []byte

	switch ptr := (content).(type) {
	case []byte:
		buf = ptr
	case string:
		buf = []byte(ptr)
	case error:
		buf = []byte(ptr.Error())
	case io.Reader:
		var status *runtime.Status

		buf, status = io2.ReadAll(io.NopCloser(ptr))
		if !status.OK() {
			return nil, "", status
		}
	case io.ReadCloser:
		var status *runtime.Status

		buf, status = io2.ReadAll(ptr)
		if !status.OK() {
			return nil, "", status
		}
	default:
		if strings.Contains(contentType, "json") {
			var status *runtime.Status

			buf, status = json.Marshal(content)
			if !status.OK() {
				return nil, "", status
			}
			return buf, contentType, status
		} else {
			return nil, "", runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
		}
	}
	if strings.Contains(contentType, "json") {
		return buf, ContentTypeJson, runtime.NewStatusOK()
	}
	return buf, http.DetectContentType(buf), runtime.NewStatusOK()
}
