package http2

import (
	"errors"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

func ReadContentFromLocation(location string) ([]byte, runtime.Status) {
	if len(location) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", errors.New("error: content-location is empty"))
	}
	u, _ := url.Parse(location)
	buf, err := io2.ReadFile(u)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", err)
	}
	return buf, runtime.NewStatusOK()
}
