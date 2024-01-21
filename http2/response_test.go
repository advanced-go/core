package http2

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/advanced-go/core/runtime"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
)

const (
	activityJsonFile = "file://[cwd]/http2test/resource/activity.json"
	activityGzipFile = "file://[cwd]/http2test/resource/activity.gz"

	testResponseHtml = "file://[cwd]/http2test/resource/test-response.html"
	jsonContentType  = "application/json"
)

type activity struct {
	ActivityID   string `json:"ActivityID"`
	ActivityType string `json:"ActivityType"`
	Agent        string `json:"Agent"`
	AgentUri     string `json:"AgentUri"`
	Assignment   string `json:"Assignment"`
	Controller   string `json:"Controller"`
	Behavior     string `json:"Behavior"`
	Description  string `json:"Description"`
}

var (
	activityJson []byte
	activityGzip []byte
	activityList []activity
)

func init() {
	var err error
	var buf []byte

	buf, err = os.ReadFile(runtime.FileName(activityJsonFile))
	if err != nil {
		fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
		return
	}
	err = json.Unmarshal(buf, &activityList)
	if err != nil {
		fmt.Printf("test: json.Unmarshal() -> [err:%v]\n", err)
		return
	}
	activityJson, err = json.Marshal(activityList)
	if err != nil {
		fmt.Printf("test: json.Mmarshal() -> [err:%v]\n", err)
		return
	}

	activityGzip, err = os.ReadFile(runtime.FileName(activityGzipFile))
	if err != nil {
		if strings.Contains(err.Error(), "open") {
			buff := new(bytes.Buffer)

			// write, flush and close
			zw := gzip.NewWriter(buff)
			cnt, err0 := zw.Write(activityJson)
			ferr := zw.Flush()
			cerr := zw.Close()
			fmt.Printf("test: gzip.Writer() -> [cnt:%v] [write-err:%v] [flush-err:%v] [close_err:%v]\n", cnt, err0, ferr, cerr)
			err = os.WriteFile(runtime.FileName(activityGzipFile), buff.Bytes(), 667)
			fmt.Printf("test: os.WriteFile(\"%v\") -> [err:%v]\n", activityGzipFile, err)

		} else {
			fmt.Printf("test: os.ReadFile() -> [err:%v]\n", err)
		}
	}
}

func ExampleWriteResponse_StatusHeaders() {
	// all nil
	rec := httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, nil, nil, nil)
	fmt.Printf("test: WriteResponse(w,nil,nil,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code
	rec = httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, nil, runtime.NewStatus(http.StatusTeapot), nil)
	fmt.Printf("test: WriteResponse(w,nil,StatusTeapot,nil) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code, headers list
	rec = httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, nil, runtime.NewStatusOK(), []Attr{{Key: ContentType, Val: ContentTypeTextHtml}, {Key: AcceptEncoding, Val: AcceptEncodingValue}})
	fmt.Printf("test: WriteResponse(w,nil,StatusOK,list) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	// status code, http.Header
	rec = httptest.NewRecorder()
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)
	h.Add(ContentEncoding, ContentEncodingGzip)
	WriteResponse[runtime.Output](rec, nil, runtime.NewStatus(http.StatusGatewayTimeout), h)
	fmt.Printf("test: WriteResponse(w,nil,StatusGatewayTimeout,http.Header) -> [status-code:%v] [header:%v]\n", rec.Result().StatusCode, rec.Result().Header)

	//Output:
	//test: WriteResponse(w,nil,nil,nil) -> [status-code:200] [header:map[]]
	//test: WriteResponse(w,nil,StatusTeapot,nil) -> [status-code:418] [header:map[]]
	//test: WriteResponse(w,nil,StatusOK,list) -> [status-code:200] [header:map[Accept-Encoding:[gzip, br] Content-Type:[text/html]]]
	//test: WriteResponse(w,nil,StatusGatewayTimeout,http.Header) -> [status-code:504] [header:map[Content-Encoding:[gzip] Content-Type:[application/json]]]

}

func ExampleWriteResponse_JSON() {
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)

	// JSON activity list
	rec := httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, activityList, runtime.StatusOK(), h)
	buf, status0 := runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:%v] [int:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	// JSON reader
	rec = httptest.NewRecorder()
	reader := bytes.NewReader(activityJson)
	WriteResponse[runtime.Output](rec, reader, runtime.StatusOK(), h)
	buf, status0 = runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: WriteResponse(w,io.Reader,OK,http.Header) -> [read-all:%v] [int:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	//Output:
	//test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:OK] [int:395] [out:395]
	//test: WriteResponse(w,io.Reader,OK,http.Header) -> [read-all:OK] [int:395] [out:395]

}

/*
func ExampleWriteResponse_StatusOK() {
	str := "text response"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusOK)
	WriteResponse[runtime.Output](w, str, status, nil)
	resp := w.Result()
	fmt.Printf("test: WriteResponse(w,%v,status) -> [status:%v] [body:%v] [header:%v]\n", str, status, w.Body.String(), resp.Header)

	//Output:
	//test: WriteResponse(w,text response,status) -> [status:OK] [body:text response] [header:map[Content-Type:[text/plain; charset=utf-8]]]

}

func ExampleWriteResponse_StatusNotOK() {
	str := "server unavailable"

	w := httptest.NewRecorder()
	status := runtime.NewStatus(http.StatusServiceUnavailable)
	WriteResponse[runtime.Output](w, str, status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	w = httptest.NewRecorder()
	status = runtime.NewStatus(http.StatusNotFound)
	WriteResponse[runtime.Output](w, []byte("not found"), status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())

	str = "operation timed out"
	w = httptest.NewRecorder()
	status = runtime.NewStatus(runtime.StatusDeadlineExceeded)
	WriteResponse[runtime.Output](w, errors.New(str), status, nil)
	fmt.Printf("test: WriteResponse(w,nil,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Header())


	//Output:
	//test: WriteResponse(w,nil,status) -> [status:503] [body:server unavailable] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:404] [body:not found] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,nil,status) -> [status:504] [body:operation timed out] [header:map[Content-Type:[text/plain; charset=utf-8]]]

}

func ExampleWriteResponse_Body() {
	w := httptest.NewRecorder()

	body := io.NopCloser(bytes.NewReader([]byte("error content")))
	WriteResponse[runtime.Output](w, body, runtime.NewStatus(http.StatusGatewayTimeout), nil)
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	body = io.NopCloser(bytes.NewReader([]byte("foo")))
	w = httptest.NewRecorder()
	WriteResponse[runtime.Output](w, body, nil,
		[]Attr{{"key", "value"}, {"key1", "value1"}, {"key2", "value2"}})
	fmt.Printf("test: WriteResponse(w,resp,status) -> [status:%v] [body:%v] [header:%v]\n", w.Code, w.Body.String(), w.Result().Header)

	//Output:
	//test: WriteResponse(w,resp,status) -> [status:504] [body:error content] [header:map[Content-Type:[text/plain; charset=utf-8]]]
	//test: WriteResponse(w,resp,status) -> [status:200] [body:foo] [header:map[Content-Type:[text/plain; charset=utf-8] Key:[value] Key1:[value1] Key2:[value2]]]

}


func ExampleWriteStatusContent() {
	r := httptest.NewRecorder()

	// No content
	writeStatusContent[runtime.Output](r, runtime.StatusOK(), "test location")
	r.Result().Header = r.Header()
	buf, status := runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Error message
	r = httptest.NewRecorder()
	writeStatusContent[runtime.Output](r, runtime.NewStatus(http.StatusInternalServerError).SetContent("error message", false), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	// Json
	d := data{Item: "test item", Count: 500}
	r = httptest.NewRecorder()
	status = runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true)
	//status.SetContent(d, true)

	writeStatusContent[runtime.Output](r, status, "test location") //runtime.NewStatus(http.StatusInternalServerError).SetContent(d, true), "test location")
	r.Result().Header = r.Header()
	buf, status = runtime.ReadAll(r.Result().Body, nil)
	fmt.Printf("test: writeStatusContent() -> %v [header:%v] [body:%v] [ReadAll:%v]\n", r.Result().StatusCode, r.Result().Header, string(buf), status)

	//Output:
	//test: writeStatusContent() -> 200 [header:map[]] [body:] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Type:[text/plain; charset=utf-8]]] [body:error message] [ReadAll:OK]
	//test: writeStatusContent() -> 200 [header:map[Content-Type:[application/json]]] [body:{"Item":"test item","Count":500}] [ReadAll:OK]

}


*/
