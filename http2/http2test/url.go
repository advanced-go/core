package http2test

import (
	"net/http"
	"net/url"
)

func UpdateUrl(newUrl string, req *http.Request) (*http.Request, error) {
	var err error
	req.URL, err = url.Parse(newUrl)
	return req, err
}
