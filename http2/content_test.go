package http2

import (
	"fmt"
)

func Example_ReadContentFromLocation() {
	buf, status := ReadContentFromLocation("file://[cwd]/http2test/resource/activity.json")
	fmt.Printf("test: ReadContentFromLocation() -> [status:%v] [content:%v]\n", status, len(buf))

	//Output:
	//test: ReadContentFromLocation() -> [status:OK] [content:525]

}
