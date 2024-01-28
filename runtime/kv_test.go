package runtime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	MSFTVariable  = "MSFT"
	MSFTAuthority = "www.bing.com"

	GOOGLVariable  = "GOOGL"
	GOOGLAuthority = "www.google.com"

	fileAttrs = "file://[cwd]/runtimetest/attrs.json"
)

func Example_KV() {
	values := []KV{{MSFTVariable, MSFTAuthority}, {GOOGLVariable, GOOGLAuthority}}

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

func ExampleStringsMap_Add() {
	smap := NewStringsMap(nil)
	key1 := "key-1"

	status := smap.Add("", "")
	fmt.Printf("test: Add(\"\") -> [status:%v]\n", status)

	status = smap.Add(key1, "value-1")
	fmt.Printf("test: Add(%v) -> [status:%v]\n", key1, status)

	//Output:
	//test: Add("") -> [status:Invalid Argument [invalid argument: key is empty]]
	//test: Add(key-1) -> [status:OK]

}

func ExampleStringsMap_Get() {
	key1 := "key-1"
	key2 := "key-2"
	h := make(http.Header)
	h.Add(key1, "value-1")
	h.Add(key2, "value-2")

	smap := NewStringsMap(h)

	val, status := smap.Get("")
	fmt.Printf("test: Get(\"\") -> [val:%v] [status:%v]\n", val, status)

	val, status = smap.Get(key1)
	fmt.Printf("test: Get(%v) -> [val:%v] [status:%v]\n", key1, val, status)

	val, status = smap.Get(key2)
	fmt.Printf("test: Get(%v) -> [val:%v] [status:%v]\n", key2, val, status)

	//Output:
	//test: Get("") -> [val:] [status:Invalid Argument [invalid argument: key does not exist: []]]
	//test: Get(key-1) -> [val:value-1] [status:OK]
	//test: Get(key-2) -> [val:value-2] [status:OK]

}
