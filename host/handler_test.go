package host

import (
	"fmt"
	"github.com/advanced-go/core/messaging"
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

func _ExamplePing() {
	uri1 := "github/advanced-go/example-domain/activity"
	r, _ := http.NewRequest("", "github/advanced-go/example-domain/activity:ping", nil)
	err := messaging.HostExchange.Add(messaging.NewMailbox(uri1, nil))
	if err != nil {
		fmt.Printf("test: processPing() -> [err:%v]\n", err)
	}
	//nid, rsc, ok := UprootUrn(r.URL.Path)
	status := messaging.Ping(nil, r.URL) //ProcessPing[runtime.Bypass](w, nid)
	fmt.Printf("test: messaging.Ping() -> [nid:%v] [nss:%v] [ok:%v] [status:%v]\n", "", "", true, status)

	//Output:
	//test: messaging.Ping() -> [nid:] [nss:] [ok:true] [status:504] [content:ping response time out: [github/advanced-go/example-domain/activity]]

}
