package log2

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
	"net/http"
	"time"
)

func Example_FmtLog() {
	start := time.Now().UTC()

	req, err := http.NewRequest("select", "https://www.google.com/search?q=test", nil)
	req.Header.Add(runtime.XRequestId, "123-456")
	fmt.Printf("test: NewRequest() -> [err:%v] %v\n", err, req)
	resp := http.Response{StatusCode: http.StatusOK}
	s := fmtLog("egress", start, time.Since(start), req, &resp, -1, "")
	fmt.Printf("test: fmtLog() -> %v\n", s)

	//Output:

}
