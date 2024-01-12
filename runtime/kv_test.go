package runtime

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	MSFTVariable  = "{MSFT}"
	MSFTAuthority = "www.bing.com"

	GOOGLVariable  = "{GOOGL}"
	GOOGLAuthority = "www.google.com"

	fileAttrs = "file://[cwd]/runtimetest/attrs.json"
)

func Example_KV() {
	values := []KV{{"MSFT", MSFTAuthority}, {"GOOGL", GOOGLAuthority}}

	buf, err := json.Marshal(values)
	fmt.Printf("test: Attr() -> [buf:%v] [err:%v]\n", string(buf), err)

	fname := FileName(fileAttrs)
	buf, err = os.ReadFile(fname)
	fmt.Printf("test: os.ReadFile(%v) -> [buf:%v] [err:%v]\n", fname, len(buf), err)
	var values2 []KV

	err = json.Unmarshal(buf, &values2)
	fmt.Printf("test: Unmarshal() -> [buf:%v] [err:%v]\n", values2, err)

	//Output:
	//test: Attr() -> [buf:[{"Key":"MSFT","Value":"www.bing.com"},{"Key":"GOOGL","Value":"www.google.com"}]] [err:<nil>]
	//test: os.ReadFile(C:\Users\markb\GitHub\core\runtime\runtimetest\attrs.json) -> [buf:124] [err:<nil>]
	//test: Unmarshal() -> [buf:[{MSFT www.bing.com} {GOOGL www.google.com}]] [err:<nil>]

}
