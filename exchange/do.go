package exchange

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/go-ai-agent/core/exchange/httptest"
	"github.com/go-ai-agent/core/runtime"
	"io"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
)

var (
	doLocation              = PkgUrl + "/do"
	Hdr        HttpExchange = Default{}
	fsys       fs.FS
)

func Do[E runtime.ErrorHandler](req *http.Request) (resp *http.Response, status *runtime.Status) {
	var e E

	if req == nil {
		return nil, e.Handle(nil, doLocation, errors.New("invalid argument : request is nil")).SetCode(runtime.StatusInvalidArgument)
	}
	var err error
	resp, err = Hdr.Do(req)
	if err != nil {
		return nil, e.Handle(req.Context(), doLocation, err).SetCode(http.StatusInternalServerError)
	}
	if resp == nil {
		return nil, e.Handle(req.Context(), doLocation, errors.New("invalid argument : response is nil")).SetCode(http.StatusInternalServerError)
	}
	return resp, runtime.NewHttpStatusCode(resp.StatusCode)
}

func DoT[E runtime.ErrorHandler, T any](req *http.Request) (resp *http.Response, t T, status *runtime.Status) {
	resp, status = Do[E](req)
	if !status.OK() {
		return nil, t, status
	}
	t, status = Deserialize[E, T](req.Context(), resp.Body)
	return
}

var http11Bytes = []byte("HTTP/1.1")
var http12Bytes = []byte("HTTP/1.2")
var http20Bytes = []byte("HTTP/2.0")

func isHttpResponseMessage(buf []byte) bool {
	if buf == nil {
		return false
	}
	l := len(buf)
	if bytes.Equal(buf[0:l], http11Bytes) {
		return true
	}
	if bytes.Equal(buf[0:l], http12Bytes) {
		return true
	}
	if bytes.Equal(buf[0:l], http20Bytes) {
		return true
	}
	return false
}

func createFileResponse(req *http.Request) (*http.Response, error) {
	path := req.URL.Path
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	buf, err := fs.ReadFile(fsys, path)
	if err != nil {
		if strings.Contains(err.Error(), "file does not exist") {
			return &http.Response{StatusCode: http.StatusNotFound}, nil
		}
		return &http.Response{StatusCode: http.StatusInternalServerError}, nil
	}
	if isHttpResponseMessage(buf) {
		return http.ReadResponse(bufio.NewReader(bytes.NewReader(buf)), req)
	} else {
		resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: &httptest.ReaderCloser{Reader: bytes.NewReader(buf), Err: nil}}
		if strings.HasSuffix(path, ".json") {
			resp.Header.Add("Content-Type", "application/json")
		} else {
			resp.Header.Add("Content-Type", "text/plain")
		}
		return resp, nil
	}
}

func createEchoResponse(req *http.Request) (*http.Response, error) {
	if req == nil {
		return nil, errors.New("invalid argument: Request is nil")
	}
	var resp = http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Request: req}
	for key, element := range req.URL.Query() {
		switch key {
		case "httpError":
			return nil, http.ErrHijacked
		case "status":
			sc, err := strconv.Atoi(element[0])
			if err == nil {
				resp.StatusCode = sc
			} else {
				resp.StatusCode = http.StatusInternalServerError
			}
		case "body":
			if len(element[0]) > 0 && resp.Body == nil {
				// Handle escaped path? See notes on the url.URL struct
				resp.Body = &httptest.ReaderCloser{Reader: strings.NewReader(element[0]), Err: nil}
			}
		case "ioError":
			// resp.StatusCode = http.StatusInternalServerError
			resp.Body = &httptest.ReaderCloser{Reader: nil, Err: io.ErrUnexpectedEOF}
		default:
			if len(element[0]) > 0 {
				resp.Header.Add(key, element[0])
			}
		}
	}
	return &resp, nil
}
