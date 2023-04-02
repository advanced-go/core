package exchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-sre/core/exchange/httptest"
	"github.com/go-sre/core/runtime"
	"net/http"
	"strings"
)

type addressV1 struct {
	Name   string
	Street string
	City   string
	ST     string
	Zip    string
}

type Zip4 struct {
	Zip   string
	Plus4 string
}

type addressV2 struct {
	Name   string
	Street string
	City   string
	State  string
	Zip    Zip4
}

func ExampleDeserialize() {
	result, status := Deserialize[runtime.DebugError2, []byte](nil, nil)
	fmt.Printf("test: Deserialize[DebugError,[]byte](nil) -> [%v] [status:%v]\n", string(result), status)

	resp := new(http.Response)
	result, status = Deserialize[runtime.DebugError2, []byte](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,[]byte](resp) -> [%v] [status:%v]\n", string(result), status)

	resp.Body = &httptest.ReaderCloser{Reader: strings.NewReader("Hello World String"), Err: nil}
	result, status = Deserialize[runtime.DebugError2, []byte](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,[]byte](resp) -> [%v] [status:%v]\n", string(result), status)

	resp.Body = &httptest.ReaderCloser{Reader: bytes.NewReader([]byte("Hello World []byte")), Err: nil}
	result2, status2 := Deserialize[runtime.DebugError2, []byte](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,[]byte](resp) -> [%v] [status:%v]\n", string(result2), status2)

	//Output:
	//[[] github.com/go-sre/core/exchange/deserialize [body is nil]]
	//test: Deserialize[DebugError,[]byte](nil) -> [] [status:Invalid Content]
	//[[] github.com/go-sre/core/exchange/deserialize [body is nil]]
	//test: Deserialize[DebugError,[]byte](resp) -> [] [status:Invalid Content]
	//test: Deserialize[DebugError,[]byte](resp) -> [Hello World String] [status:OK]
	//test: Deserialize[DebugError,[]byte](resp) -> [Hello World []byte] [status:OK]

}

func ExampleDeserialize_Decode() {
	addrV1 := addressV1{
		Name:   "Bob Smith",
		Street: "123 Oak Avenue",
		City:   "New Orleans",
		ST:     "LA",
		Zip:    "12345",
	}
	bufV1, _ := json.Marshal(&addrV1)

	resp := new(http.Response)
	resp.Body = &httptest.ReaderCloser{Reader: bytes.NewReader(bufV1), Err: nil}

	result, status := Deserialize[runtime.DebugError2, addressV1](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,addressV1](resp) -> [%v] [status:%v]\n", result, status)

	addrV2 := addressV2{
		Name:   "Bob Smith",
		Street: "123 Oak Avenue",
		City:   "New Orleans",
		State:  "Louisiana",
		Zip:    Zip4{Zip: "12345", Plus4: "1234"},
	}
	bufV2, _ := json.Marshal(&addrV2)
	resp = new(http.Response)
	resp.Body = &httptest.ReaderCloser{Reader: bytes.NewReader(bufV2), Err: nil}

	result2, status2 := Deserialize[runtime.DebugError2, addressV2](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,addressV2](resp) -> [%v] [status:%v]\n", result2, status2)

	resp = new(http.Response)
	resp.Body = &httptest.ReaderCloser{Reader: bytes.NewReader(bufV2), Err: nil}

	result3, status3 := Deserialize[runtime.DebugError2, addressV1](nil, resp.Body)
	fmt.Printf("test: Deserialize[DebugError,addressV1](resp) -> [%v] [status:%v]\n", result3, status3)

	//Output:
	//test: Deserialize[DebugError,addressV1](resp) -> [{Bob Smith 123 Oak Avenue New Orleans LA 12345}] [status:OK]
	//test: Deserialize[DebugError,addressV2](resp) -> [{Bob Smith 123 Oak Avenue New Orleans Louisiana {12345 1234}}] [status:OK]
	//[[] github.com/go-sre/core/exchange/deserialize [json: cannot unmarshal object into Go struct field addressV1.Zip of type string]]
	//test: Deserialize[DebugError,addressV1](resp) -> [{Bob Smith 123 Oak Avenue New Orleans  }] [status:Json Decode Failure]
	
}
