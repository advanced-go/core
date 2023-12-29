package uri

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
)

const (
	CwdVariable = "[cwd]"
	statusToken = "status"
	JsonExt     = ".json"
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

func IsJson(uri string) bool {
	return strings.HasSuffix(uri, JsonExt)
}

func IsFileScheme(u any) bool {
	if u == nil {
		return false
	}
	if s, ok := u.(string); ok {
		return strings.HasPrefix(s, FileScheme)
	}
	if u2, ok := u.(*url.URL); ok {
		return u2.Scheme == FileScheme
	}
	return false
}

func FileName(uri any) string {
	if uri == nil {
		return "error: URL is nil"
	}
	if s, ok := uri.(string); ok {
		if len(s) == 0 {
			return "error: URL is empty"
		}
		return fileName(ParseRaw(s))
	}
	if u, ok := uri.(*url.URL); ok {
		return fileName(u)
	}
	return fmt.Sprintf("error: invalid URL type: %v", reflect.TypeOf(uri))
}

func fileName(u *url.URL) string {
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
