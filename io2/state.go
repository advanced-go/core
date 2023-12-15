package io2

import (
	"encoding/json"
	"errors"
	"github.com/advanced-go/core/runtime"
	"net/url"
)

const (
	readStateLoc   = PkgPath + ":ReadState"
	readResultsLoc = PkgPath + ":ReadResults"
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

func ReadResults[T any](urls []string) (t T, status runtime.Status) {
	switch len(urls) {
	case 0:
		return t, runtime.StatusOK()
	case 1:
		if len(urls[0]) == 0 {
			return t, runtime.StatusOK()
		}
		return ReadState[T](urls[0])
	default:
	}
	if len(urls[0]) > 0 {
		t, status = ReadState[T](urls[0])
		if !status.OK() {
			return t, status
		}
	}
	return t, ReadStatus(urls[1])
}
