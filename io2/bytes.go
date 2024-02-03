package io2

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
	bytesLoc  = runtime.PkgPath + ":Bytes"
	jsonToken = "json"
)

// Bytes - convert content to []byte, checking for JSON content
func Bytes(content any, contentType string) ([]byte, *runtime.Status) {
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

		buf, status = ReadAll(ptr, nil)
		if !status.OK() {
			return nil, status
		}
	case io.ReadCloser:
		var status *runtime.Status

		buf, status = ReadAll(ptr, nil)
		_ = ptr.Close()
		if !status.OK() {
			return nil, status
		}
	default:
		if strings.Contains(contentType, jsonToken) {
			var err error

			buf, err = json.Marshal(content)
			if err != nil {
				status := runtime.NewStatusError(runtime.StatusJsonEncodeError, bytesLoc, err)
				if !status.OK() {
					return nil, status
				}
			}
			return buf, runtime.StatusOK()
		} else {
			return nil, runtime.NewStatusError(http.StatusInternalServerError, bytesLoc, errors.New(fmt.Sprintf("error: content type is invalid: %v", reflect.TypeOf(ptr))))
		}
	}
	return buf, runtime.StatusOK()
}

/*
	case *http.Response:
		var status Status

		buf, status = ReadAll(ptr.Body, nil)
		_ = ptr.Body.Close()
		if !status.OK() {
			return nil, status
		}

*/