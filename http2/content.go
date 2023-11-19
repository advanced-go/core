package http2

import (
	"errors"
	"github.com/advanced-go/core/io2"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/url"
)

func ReadContentFromLocation(h http.Header) ([]byte, runtime.Status) {
	content := h.Get(ContentLocation)
	if len(content) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", errors.New("error: content-location is empty"))
	}
	u, _ := url.Parse(content)
	buf, err := io2.ReadFile(u)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadContentFromLocation()", err)
	}
	return buf, runtime.NewStatusOK()
}
