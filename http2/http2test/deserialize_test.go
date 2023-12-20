package http2test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	result, status := Deserialize[[]byte](nil)
	fmt.Printf("test: Deserialize[[]byte](nil) -> [%v] [status:%v]\n", string(result), status)

	resp := new(http.Response)
	result, status = Deserialize[[]byte](resp.Body)
	fmt.Printf("test: Deserialize[[]byte](resp) -> [%v] [status:%v]\n", string(result), status)

	resp.Body = io.NopCloser(strings.NewReader("Hello World String"))
	result, status = Deserialize[[]byte](resp.Body)
	fmt.Printf("test: Deserialize[[]byte](resp) -> [%v] [status:%v]\n", string(result), status)

	resp.Body = io.NopCloser(bytes.NewReader([]byte("Hello World []byte")))
	result2, status2 := Deserialize[[]byte](resp.Body)
	fmt.Printf("test: Deserialize[[]byte](resp) -> [%v] [status:%v]\n", string(result2), status2)

	//Output:
	//test: Deserialize[[]byte](nil) -> [] [status:Invalid Content [body is nil]]
	//test: Deserialize[[]byte](resp) -> [] [status:Invalid Content [body is nil]]
	//test: Deserialize[[]byte](resp) -> [Hello World String] [status:OK]
	//test: Deserialize[[]byte](resp) -> [Hello World []byte] [status:OK]

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
	resp.Body = io.NopCloser(bytes.NewReader(bufV1))

	result, status := Deserialize[addressV1](resp.Body)
	fmt.Printf("test: Deserialize[addressV1](resp) -> [%v] [status:%v]\n", result, status)

	addrV2 := addressV2{
		Name:   "Bob Smith",
		Street: "123 Oak Avenue",
		City:   "New Orleans",
		State:  "Louisiana",
		Zip:    Zip4{Zip: "12345", Plus4: "1234"},
	}
	bufV2, _ := json.Marshal(&addrV2)
	resp = new(http.Response)
	resp.Body = io.NopCloser(bytes.NewReader(bufV2))

	result2, status2 := Deserialize[addressV2](resp.Body)
	fmt.Printf("test: Deserialize[addressV2](resp) -> [%v] [status:%v]\n", result2, status2)

	resp = new(http.Response)
	resp.Body = io.NopCloser(bytes.NewReader(bufV2))

	result3, status3 := Deserialize[addressV1](resp.Body)
	fmt.Printf("test: Deserialize[addressV1](resp) -> [%v] [status:%v]\n", result3, status3)

	//Output:
	//test: Deserialize[addressV1](resp) -> [{Bob Smith 123 Oak Avenue New Orleans LA 12345}] [status:OK]
	//test: Deserialize[addressV2](resp) -> [{Bob Smith 123 Oak Avenue New Orleans Louisiana {12345 1234}}] [status:OK]
	//test: Deserialize[addressV1](resp) -> [{Bob Smith 123 Oak Avenue New Orleans  }] [status:Json Decode Failure [json: cannot unmarshal object into Go struct field addressV1.Zip of type string]]

}
