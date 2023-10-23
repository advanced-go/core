package httptest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"testing"
)

type Args struct {
	Item string
	Got  string
	Want string
	Err  error
}

type Fail func(a []Args, t *testing.T)

func Headers(got *http.Response, want *http.Response, t *testing.T, fail Fail, names ...string) {
	var failures []Args

	for _, name := range names {
		wantVal := want.Header.Get(name)
		if wantVal == "" {
			fail([]Args{{Item: name, Got: "", Want: "", Err: errors.New(fmt.Sprintf("want header [%v] is missing or empty", name))}}, t)
			return
		}
		gotVal := got.Header.Get(name)
		if wantVal != gotVal {
			failures = append(failures, Args{Item: name, Got: gotVal, Want: wantVal, Err: nil})
		}
	}
	if failures != nil {
		fail(failures, t)
	}
}

func Content[T any](got *http.Response, want *http.Response, t *testing.T, fail Fail, test func(got T, want T, t *testing.T)) {
	// validate content type
	valid, ct := validateContentType(got, want, t, fail)
	if !valid {
		return
	}

	// validate body IO
	wantBytes, status := exchange.ReadAll[runtime.BypassError](want.Body)
	if status.IsErrors() {
		fail([]Args{{Item: "want.Body", Got: "", Want: "", Err: status.Errors()[0]}}, t)
		return
	}
	gotBytes, status1 := exchange.ReadAll[runtime.BypassError](got.Body)
	if status1.IsErrors() {
		fail([]Args{{Item: "got.Body", Got: "", Want: "", Err: status1.Errors()[0]}}, t)
		return
	}

	// validate content length
	if len(gotBytes) != len(wantBytes) {
		fail([]Args{{Item: "Content-Length", Got: fmt.Sprintf("%v", len(gotBytes)), Want: fmt.Sprintf("%v", len(wantBytes))}}, t)
		return
	}

	// if no content or content type is no application/json, return
	if len(wantBytes) == 0 || ct != runtime.ContentTypeJson {
		return
	}

	// unmarshal
	var gotT T
	var wantT T

	err := json.Unmarshal(wantBytes, &wantT)
	if err != nil {
		fail([]Args{{Item: "want.Unmarshal()", Got: "", Want: "", Err: err}}, t)
		return
	}
	err = json.Unmarshal(gotBytes, &gotT)
	if err != nil {
		fail([]Args{{Item: "got.Unmarshal()", Got: "", Want: "", Err: err}}, t)
		return
	}

	// test
	test(gotT, wantT, t)
}

func ReadHttp(basePath, reqName, respName string, t *testing.T, fail Fail) (*http.Request, *http.Response) {
	path := basePath + reqName
	req, err := exchange.ReadRequest(runtime.ParseRaw(path))
	if err != nil {
		fail([]Args{{Item: "ReadRequest", Got: "", Want: "", Err: errors.New(fmt.Sprintf("ReadRequest(%v) failed : %v", path, err))}}, t)
		return nil, nil
	}
	path = basePath + respName
	resp, err1 := exchange.ReadResponse(runtime.ParseRaw(path))
	if err1 != nil {
		fail([]Args{{Item: "ReadResponse", Got: "", Want: "", Err: errors.New(fmt.Sprintf("ReadResponse(%v) failed : %v", path, err1))}}, t)
		return nil, nil
	}
	return req, resp
}

func validateContentType(got *http.Response, want *http.Response, t *testing.T, fail Fail) (bool, string) {
	ct := want.Header.Get(runtime.ContentType)
	if ct == "" {
		fail([]Args{{Item: runtime.ContentType, Got: "", Want: "", Err: errors.New("want Response header Content-Type is empty")}}, t)
		return false, ct
	}
	gotCt := got.Header.Get(runtime.ContentType)
	if gotCt != ct {
		fail([]Args{{Item: runtime.ContentType, Got: gotCt, Want: ct, Err: nil}}, t)
		return false, ct
	}
	return true, ct
}
