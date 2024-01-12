package runtime

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
	jsonExt     = ".json"
	fileScheme  = "file"
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

// IsJsonURL - does the URI have a .json extension
func IsJsonURL(uri string) bool {
	return strings.HasSuffix(uri, jsonExt)
}

// IsStatusURL - determine if the file name of the URL contains the text 'status'
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

// FileName - return the OS correct file name from a URI
func FileName(uri any) string {
	if uri == nil {
		return "error: URL is nil"
	}
	if s, ok := uri.(string); ok {
		if len(s) == 0 {
			return "error: URL is empty"
		}
		u, _ := url.Parse(s)
		return fileName(u)
	}
	if u, ok := uri.(*url.URL); ok {
		return fileName(u)
	}
	return fmt.Sprintf("error: invalid URL type: %v", reflect.TypeOf(uri))
}

func fileName(u *url.URL) string {
	if u == nil || u.Scheme != fileScheme {
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

// isFileScheme - determine if a string, or URL uses the file scheme
/*
func isFileScheme(u any) bool {
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
*/
