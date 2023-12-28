package io2

import (
	"github.com/advanced-go/core/runtime"
	"io"
)

const (
	PkgPath = "github.com/advanced-go/http2/io2"
)

// ReadAll - read the body with a runtime.Status
func ReadAll(body io.ReadCloser) ([]byte, runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, PkgPath+":ReadAll", err)
	}
	return buf, runtime.StatusOK()
}
