package exchange

import (
	"fmt"
	"net/http"
	"os"
)

func ExampleLookupRequest() {
	uri := "http://localhost:8080/base-path/resource?first=false"
	req, _ := http.NewRequest("", uri, nil)

	name := ""

	s, err := LookupRequest("", nil)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("test", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("method", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("scheme", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("host", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("path", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	s, err = LookupRequest("query", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	os.Setenv("RUNTIME", "DEV")
	s, err = LookupRequest("$RUNTIME", req)
	fmt.Printf("test: LookupRequest(%v) -> [err:%v] [%v]\n", name, err, s)

	//Output:
	//test: LookupRequest() -> [err:invalid argument: Request is nil] []
	//test: LookupRequest() -> [err:invalid argument : LookupEnv() template variable is invalid: test] []
	//test: LookupRequest() -> [err:<nil>] [GET]
	//test: LookupRequest() -> [err:<nil>] [http]
	//test: LookupRequest() -> [err:<nil>] [localhost:8080]
	//test: LookupRequest() -> [err:<nil>] [/base-path/host]
	//test: LookupRequest() -> [err:<nil>] [first=false]
	//test: LookupRequest() -> [err:<nil>] [DEV]

}
