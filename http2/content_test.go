package http2

import (
	"fmt"
	"net/http"
)

func Example_ReadContentFromLocation() {
	h := make(http.Header)
	h.Add(ContentLocation, "file://[cwd]/http2test/resource/activity.json")
	buf, status := ReadContentFromLocation(h)

	fmt.Printf("test: ReadContentFromLocation() -> [status:%v] [content:%v]\n", status, len(buf))

	//Output:
	//test: ReadContentFromLocation() -> [status:OK] [content:525]

}
