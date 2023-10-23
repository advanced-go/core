package httpx

import (
	"fmt"
	"github.com/go-ai-agent/core/runtime"
)

func Example_ReadContent() {
	s := "file://[cwd]/httpxtest/resource/http/get-request.txt"
	buf, err := ReadFile(runtime.ParseRaw(s))
	if err != nil {
		fmt.Printf("test: ReadFile(%v) -> [err:%v]\n", s, err)

	} else {
		bytes, err1 := ReadContent(buf)
		fmt.Printf("test: ReadContent() -> [err:%v] [bytes:%v]\n", err1, bytes.Len())
	}

	//Output:

}
