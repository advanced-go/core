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
	//test: WriteResponse(w,nil,StatusOK,list) -> [status-code:200] [header:map[Accept-Encoding:[gzip, deflate, br] Content-Type:[text/html]]]
	//test: WriteResponse(w,nil,StatusGatewayTimeout,http.Header) -> [status-code:504] [header:map[Content-Encoding:[gzip] Content-Type:[application/json]]]

}

func ExampleWriteResponse_JSON() {
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)

	// JSON activity list
	rec := httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, activityList, runtime.StatusOK(), h)
	buf, status0 := runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:%v] [in:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	// JSON reader
	rec = httptest.NewRecorder()
	reader := bytes.NewReader(activityJson)
	WriteResponse[runtime.Output](rec, reader, runtime.StatusOK(), h)
	buf, status0 = runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: WriteResponse(w,io.Reader,OK,http.Header) -> [read-all:%v] [in:%v] [out:%v]\n", status0, len(activityJson), len(buf))

	//Output:
	//test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:OK] [in:395] [out:395]
	//test: WriteResponse(w,io.Reader,OK,http.Header) -> [read-all:OK] [in:395] [out:395]

}

func ExampleWriteResponse_Encoding() {
	h := make(http.Header)
	h.Add(ContentType, ContentTypeJson)
	h.Add(AcceptEncoding, AcceptEncodingValue)

	// JSON activity list
	rec := httptest.NewRecorder()
	WriteResponse[runtime.Output](rec, activityList, runtime.StatusOK(), h)
	buf, status0 := runtime.ReadAll(rec.Result().Body, nil)
	fmt.Printf("test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:%v] [buf:%v][header:%v]\n", status0, http.DetectContentType(buf), rec.Result().Header)

	//Output:
	//test: WriteResponse(w,[]activity,OK,http.Header) -> [read-all:OK] [buf:application/x-gzip][header:map[Content-Encoding:[gzip] Content-Type:[application/json]]]

}
