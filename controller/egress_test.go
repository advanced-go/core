package controller

import (
	"fmt"
	"github.com/advanced-go/core/runtime"
	"io"
	"net/http"
	"time"
)

func testDo(req *http.Request) (*http.Response, *runtime.Status) {
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		return resp, runtime.NewStatus(resp.StatusCode)
	}
	resp = &http.Response{StatusCode: runtime.StatusDeadlineExceeded}
	return resp, runtime.NewStatusError(runtime.StatusDeadlineExceeded, err)
}

func ExampleDoEgress() {
	uri := "https://www.google.com/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := doEgress(time.Second*5, testDo, req)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: ExampleDoEgress_OK -> [status-code:%v] [status:%v] [content:%v]\n", resp.StatusCode, status, len(buf) > 0)

	/*
		buf = nil
		resp, status = doEgress(time.Millisecond*5, testDo, req)
		if resp.Body != nil {
			buf, _ = io.ReadAll(resp.Body)
		}
		fmt.Printf("test: ExampleDoEgress_Timeout -> [status-code:%v] [status:%v] [content:%v]\n", resp.StatusCode, status, len(buf) > 0)
	*/

	resp, status = doEgress(time.Millisecond*500, func(r *http.Request) (*http.Response, *runtime.Status) {
		time.Sleep(time.Second * 2)
		return testDo(r)
	}, req)
	time.Sleep(time.Second * 3)
	fmt.Printf("test: ExampleDoEgress_Recover -> [status-code:%v] [status:%v] [content:%v]\n", resp.StatusCode, status, false)

	//Output:
	//test: ExampleDoEgress_OK -> [status-code:200] [status:OK] [content:true]
	//test: recovered in controller.doEgress() : send on closed channel
	//test: ExampleDoEgress_Recover -> [status-code:504] [status:Deadline Exceeded [context deadline exceeded]] [content:false]

}
