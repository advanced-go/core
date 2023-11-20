package io2

import (
	"errors"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/url"
	"os"
	"strings"
)

var (
	basePath = ""
	win      = false
)

// init - set the base path and windows flag
func init() {
	cwd, err := os.Getwd()
	if err != nil {
		basePath = err.Error()
	}
	if os.IsPathSeparator(uint8(92)) {
		win = true
	}
	basePath = cwd
}

// ReadFile - read a file from the given URL template
func ReadFile(u *url.URL) ([]byte, error) {
	if u == nil {
		return nil, errors.New("error: Url is nil")
	}
	if u.Scheme != "file" {
		return nil, errors.New(fmt.Sprintf("error: Scheme is not valid [%v]", u.Scheme))
	}
	name := createFname(u)
	return os.ReadFile(name)
}

func createFname(u *url.URL) string {
	name := basePath
	if u.Host == "[cwd]" {
		name += u.Path
	} else {
		name = u.Path[1:]
	}
	if win {
		name = strings.ReplaceAll(name, "/", "\\")
	}
	return name
}

// ReadFileFromPath - read a file given a templated path
func ReadFileFromPath(path string) ([]byte, runtime.Status) {
	if len(path) == 0 {
		return nil, runtime.NewStatusError(runtime.StatusInvalidArgument, "ReadFileFromPath()", errors.New("error: path is empty"))
	}
	u, _ := url.Parse(path)
	buf, err := ReadFile(u)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusInvalidContent, "ReadFileFromPath()", err)
	}
	return buf, runtime.NewStatusOK()
}
