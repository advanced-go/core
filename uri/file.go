package uri

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

const (
	CwdVariable = "[cwd]"
	statusToken = "status"
	//statusSegment     = "/status/"
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

func IsFileScheme(u *url.URL) bool {
	if u == nil {
		return false
	}
	return u.Scheme == FileScheme
}

func FileName(u *url.URL) string {
	if u == nil {
		return "error: URL is nil"
	}
	if !IsFileScheme(u) {
		return fmt.Sprintf("error: scheme is invalid [%v]", u.Scheme)
	}
	name := basePath
	if u.Host == CwdVariable {
		name += u.Path
	} else {
		name = u.Path[1:]
	}
	if win {
		name = strings.ReplaceAll(name, "/", "\\")
	}
	return name
}

func IsStatusURL(url string) bool {
	if len(url) == 0 {
		return false
	}
	i := strings.LastIndex(url, statusToken)
	if i == -1 {
		return false
	}
	return strings.LastIndex(url, "/") < i
}
