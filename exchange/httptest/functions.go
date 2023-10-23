package httptest

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-ai-agent/core/exchange"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
)

type Args struct {
	Item string
	Got  string
	Want string
	Err  error
}

func ReadHttp(basePath, reqName, respName string) ([]Args, *http.Request, *http.Response) {
	path := basePath + reqName
	req, err := exchange.ReadRequest(runtime.ParseRaw(path))
	if err != nil {
		return []Args{{Item: "ReadRequest()", Got: "", Want: "", Err: err}}, nil, nil
	}
	path = basePath + respName
	resp, err1 := exchange.ReadResponse(runtime.ParseRaw(path))
	if err1 != nil {
		return []Args{{Item: "ReadResponse()", Got: "", Want: "", Err: err1}}, nil, nil
	}
	return nil, req, resp
}

func Headers(got *http.Response, want *http.Response, names ...string) (failures []Args) {
	for _, name := range names {
		wantVal := want.Header.Get(name)
		if wantVal == "" {
			return []Args{{Item: name, Got: "", Want: "", Err: errors.New(fmt.Sprintf("want header [%v] is missing or empty", name))}}
		}
		gotVal := got.Header.Get(name)
		if wantVal != gotVal {
			failures = append(failures, Args{Item: name, Got: gotVal, Want: wantVal, Err: nil})
		}
	}
	return failures
}

func Content[T any](got *http.Response, want *http.Response) (failures []Args, gotT T, wantT T) {
	// validate content type
	fails, ct := validateContentType(got, want)
	if fails != nil {
		failures = fails
		return
	}

	// validate body IO
	wantBytes, status := exchange.ReadAll[runtime.BypassError](want.Body)
	if status.IsErrors() {
		failures = []Args{{Item: "want.Body", Got: "", Want: "", Err: status.Errors()[0]}}
		return
	}
	gotBytes, status1 := exchange.ReadAll[runtime.BypassError](got.Body)
	if status1.IsErrors() {
		failures = []Args{{Item: "got.Body", Got: "", Want: "", Err: status1.Errors()[0]}}
		return
	}

	// validate content length
	if len(gotBytes) != len(wantBytes) {
		failures = []Args{{Item: "Content-Length", Got: fmt.Sprintf("%v", len(gotBytes)), Want: fmt.Sprintf("%v", len(wantBytes))}}
		return
	}

	// if no content or content type is no application/json, return
	if len(wantBytes) == 0 || ct != runtime.ContentTypeJson {
		return
	}

	// unmarshal
	err := json.Unmarshal(wantBytes, &wantT)
	if err != nil {
		failures = []Args{{Item: "want.Unmarshal()", Got: "", Want: "", Err: err}}
		return
	}
	err = json.Unmarshal(gotBytes, &gotT)
	if err != nil {
		failures = []Args{{Item: "got.Unmarshal()", Got: "", Want: "", Err: err}}
	}
	return
}

func validateContentType(got *http.Response, want *http.Response) (failures []Args, ct string) {
	ct = want.Header.Get(runtime.ContentType)
	if ct == "" {
		return []Args{{Item: runtime.ContentType, Got: "", Want: "", Err: errors.New("want Response header Content-Type is empty")}}, ct
	}
	gotCt := got.Header.Get(runtime.ContentType)
	if gotCt != ct {
		return []Args{{Item: runtime.ContentType, Got: gotCt, Want: ct, Err: nil}}, ct
	}
	return nil, ct
}
