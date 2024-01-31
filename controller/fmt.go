package controller

import (
	"net/http"
	"strings"
	"time"
)

// Milliseconds - convert time.Duration to milliseconds
func Milliseconds(duration time.Duration) int {
	return int(duration / time.Duration(1e6))
}

// CreateUrlHostPath - create the URL, host and path
func CreateUrlHostPath(req *http.Request) (url string, host string, path string) {
	host = req.Host
	if len(host) == 0 {
		host = req.URL.Host
	}
	url = req.URL.String()
	if len(host) == 0 {
		//url = "urn:" + url
	} else {
		if len(req.URL.Scheme) == 0 {
			url = "http://" + host + req.URL.Path
		}
	}
	path = req.URL.Path
	i := strings.Index(path, ":")
	if i >= 0 {
		path = path[i+1:]
	}
	return
}

// FmtJsonString - Json format a string value
func FmtJsonString(value string) string {
	if len(value) == 0 {
		return "null"
	}
	return "\"" + value + "\""
}
