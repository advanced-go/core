package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

const (
	readStateLoc = PkgPath + ":ReadState"
)

func ReadState[T any](uri string) (t T, status runtime.Status) {
	if len(uri) == 0 {
		return t, runtime.NewStatusError(runtime.StatusInvalidArgument, readStateLoc, errors.New("Uir is empty"))
	}
	u, err := url.Parse(uri)
	if err != nil {
		return t, runtime.NewStatusError(runtime.StatusInvalidArgument, readStateLoc, err)
	}
	buf, err1 := ReadFile(u)
	if err != nil {
		return t, runtime.NewStatusError(runtime.StatusIOError, readStateLoc, err1)
	}
	err = json.Unmarshal(buf, &t)
	if err != nil {
		return t, runtime.NewStatusError(runtime.StatusJsonDecodeError, readStateLoc, err)
	}
	return t, runtime.StatusOK()
}
