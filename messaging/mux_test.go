package messaging

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"net/http/httptest"
)

func appHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusGatewayTimeout)
}

func Example_HttpHandler() {
	pattern := "github.com/advanced-go/example-domain/activity"
	r, _ := http.NewRequest("PUT", "http://localhost:8080/github.com/advanced-go/example-domain/activity:entry", nil)

	Handle(pattern, appHttpHandler)

	rec := httptest.NewRecorder()

	HttpHandler(rec, r)

	fmt.Printf("test: HttpHandler() -> %v\n", rec.Result().StatusCode)

	//Output:
	//test: HttpHandler() -> 504

}

func Example_ProcessPing() {
	uri1 := "github.com/advanced-go/example-domain/activity"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "github.com/advanced-go/example-domain/activity:ping", nil)
	status := HostExchange.Add(NewMailbox(uri1, nil))
	if !status.OK() {
		fmt.Printf("test: processPing() -> [status:%v]\n", status)
	}
	nid, rsc, ok := uprootUrn(r.URL.Path)
	ProcessPing[runtime.Bypass](w, nid)
	buf, status1 := runtime.NewBytes(w.Result())
	if !status1.OK() {
		fmt.Printf("test: ReadAll() -> [status:%v]\n", status1)
	}
	fmt.Printf("test: processPing() -> [nid:%v] [nss:%v] [ok:%v] [status:%v] [content:%v]\n", nid, rsc, ok, w.Result().StatusCode, string(buf))

	//Output:
	//test: processPing() -> [nid:github.com/advanced-go/example-domain/activity] [nss:ping] [ok:true] [status:504] [content:]

}

func readAll(body io.ReadCloser) ([]byte, runtime.Status) {
	if body == nil {
		return nil, runtime.StatusOK()
	}
	defer body.Close()
	buf, err := io.ReadAll(body)
	if err != nil {
		return nil, runtime.NewStatusError(runtime.StatusIOError, PkgPath+":ReadAll", err)
	}
	return buf, runtime.StatusOK()
}
