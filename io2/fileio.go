package io2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"github.com/advanced-go/core/uri"
	"net/url"
	"os"
)

// ReadFile - read a file from the given URL template
func ReadFile_UNUSED(u *url.URL) ([]byte, error) {
	if u == nil {
		return nil, errors.New("error: Url is nil")
	}
	if u.Scheme != "file" {
		return nil, errors.New(fmt.Sprintf("error: Scheme is not valid [%v]", u.Scheme))
	}
	name := uri.FileName(u)
	return os.ReadFile(name)
}

// ReadFileFromPath - read a file given a templated path
func ReadFileFromPath_UNUSED(path string) ([]byte, runtime.Status) {
	if len(path) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, "ReadFileFromPath()", errors.New("error: path is empty"))
	}
	u, _ := url.Parse(path)
	buf, err := ReadFile_UNUSED(u)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadFileFromPath()", err)
	}
	return buf, runtime.StatusOK()
}
