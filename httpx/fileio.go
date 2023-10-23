package httpx

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var (
	basePath = ""
	win      = false
)

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
