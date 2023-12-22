package uri

import (
	"fmt"
	"net/http"
	"os"
)

func ExampleLookupRequest() {
	uri := "http://localhost:8080/base-path/resource?first=false"
	req, _ := http.NewRequest("", uri, nil)

	name := ""

	s, err := LookupVariable("", nil)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = "test"
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	//s, err = LookupVariable("method", req.URL)
	//fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = SchemeName
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = HostName
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = PathName
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = QueryName
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	name = "$RUNTIME"
	os.Setenv("RUNTIME", "DEV")
	s, err = LookupVariable(name, req.URL)
	fmt.Printf("test: LookupVariable(%v) -> [err:%v] [%v]\n", name, err, s)

	//Output:
	//test: LookupVariable() -> [err:invalid argument: Url is nil] []
	//test: LookupVariable(test) -> [err:invalid argument : LookupEnv() template variable is invalid: test] []
	//test: LookupVariable(scheme) -> [err:<nil>] [http]
	//test: LookupVariable(host) -> [err:<nil>] [localhost:8080]
	//test: LookupVariable(path) -> [err:<nil>] [/base-path/resource]
	//test: LookupVariable(query) -> [err:<nil>] [first=false]
	//test: LookupVariable($RUNTIME) -> [err:<nil>] [DEV]

}
