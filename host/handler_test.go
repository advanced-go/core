package host

import (
	"fmt"
	"github.com/advanced-go/core/messaging"
	"io"
	"net/http"
	"net/http/httptest"
)

func appHttpHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTeapot)
}

func Example_HttpHandler() {
	pattern := "github/advanced-go/example-domain/activity"
	r, _ := http.NewRequest("PUT", "http://localhost:8080/github/advanced-go/example-domain/activity:entry", nil)
	RegisterHandler(pattern, appHttpHandler)

	rec := httptest.NewRecorder()
	HttpHandler(rec, r)

	fmt.Printf("test: HttpHandler() -> %v\n", rec.Result().StatusCode)

	//Output:
	//test: HttpHandler() -> 418

}

func _Example_ProcessPing() {
	uri1 := "github/advanced-go/example-domain/activity"
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("", "github/advanced-go/example-domain/activity:ping", nil)
	err := messaging.HostExchange.Add(messaging.NewMailbox(uri1, nil))
	if err != nil {
		fmt.Printf("test: processPing() -> [err:%v]\n", err)
	}
	nid, rsc, ok := uprootUrn(r.URL.Path)
	//ProcessPing[runtime.Bypass](w, nid)
	buf, err0 := io.ReadAll(w.Result().Body)
	if err0 != nil {
		fmt.Printf("test: ReadAll() -> [status:%v]\n", err0)
	}
	fmt.Printf("test: processPing() -> [nid:%v] [nss:%v] [ok:%v] [status:%v] [content:%v]\n", nid, rsc, ok, w.Result().StatusCode, string(buf))

	//Output:
	//test: processPing() -> [nid:github/advanced-go/example-domain/activity] [nss:ping] [ok:true] [status:504] [content:ping response time out: [github/advanced-go/example-domain/activity]]

}
