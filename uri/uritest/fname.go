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
	JsonExt     = ".json"
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

// IsFileScheme - determine if a string, or URL uses the file scheme
func IsFileScheme(u any) bool {
	if u == nil {
		return false
	}
	if s, ok := u.(string); ok {
		return strings.HasPrefix(s, fileScheme)
	}
	if u2, ok := u.(*url.URL); ok {
		return u2.Scheme == fileScheme
	}
	return false
}

// FileName2 - return the OS correct file name from a URI
func FileName2(uri any) string {
	if uri == nil {
		return "error: URL is nil"
	}
	if s, ok := uri.(string); ok {
		if len(s) == 0 {
			return "error: URL is empty"
		}
		return fileName2(parseRaw(s))
	}
	if u, ok := uri.(*url.URL); ok {
		return fileName2(u)
	}
	return fmt.Sprintf("error: invalid URL type: %v", reflect.TypeOf(uri))
}

func fileName2(u *url.URL) string {
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

// ParseRaw - parse a raw Uri without error
func parseRaw(rawUri string) *url.URL {
	u, _ := url.Parse(rawUri)
	return u
}
