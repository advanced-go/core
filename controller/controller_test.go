package controller

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func ExampleDo_Internal() {
	ctrl := NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", 0, httpCall), nil)
	uri := "/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := ctrl.Do(nil, req)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", time.Millisecond*5, httpCall), nil)
	resp, status = ctrl.Do(nil, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_5ms() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: httpCall() -> [content:134585] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: Do_0s() -> [status-code:200] [status:OK] [buff:true]
	//test: httpCall() -> [content:0] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>] [write-err:<nil>]
	//test: Do_5ms() -> [status-code:504] [status:Timeout] [buff:true]

}

func ExampleDo_Internal_Deadline() {
	ctrl := NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", 0, httpCall), nil)
	uri := "/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := ctrl.Do(nil, req)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", time.Millisecond*5, httpCall), nil)
	resp, status = ctrl.Do(nil, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_5ms() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: httpCall() -> [content:134585] [do-err:<nil>] [read-err:<nil>] [write-err:<nil>]
	//test: Do_0s() -> [status-code:200] [status:OK] [buff:true]
	//test: httpCall() -> [content:0] [do-err:Get "https://www.google.com/search?q=golang": context deadline exceeded] [read-err:<nil>] [write-err:<nil>]
	//test: Do_5ms() -> [status-code:504] [status:Timeout] [buff:true]

}

func ExampleDo_Egress() {
	ctrl := NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", 0, nil), nil)
	uri := "/search?q=golang"
	req, _ := http.NewRequest(http.MethodGet, uri, nil)

	resp, status := ctrl.Do(testDo, req)
	buf, _ := io.ReadAll(resp.Body)
	fmt.Printf("test: Do_0s() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	ctrl = NewController("google-search", NewPrimaryResource("https://www.google.com", "/health/liveness", time.Millisecond*5, nil), nil)
	resp, status = ctrl.Do(testDo, req)
	if status.OK() {
		buf, _ = io.ReadAll(resp.Body)
	}
	fmt.Printf("test: Do_5ms() -> [status-code:%v] [status:%v] [buff:%v]\n", resp.StatusCode, status, len(buf) > 0)

	//Output:
	//test: Do_0s() -> [status-code:200] [status:OK] [buff:true]
	//test: Do_5ms() -> [status-code:504] [status:Deadline Exceeded [context deadline exceeded]] [buff:true]

}
