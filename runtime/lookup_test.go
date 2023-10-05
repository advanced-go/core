package runtime

import (
	"fmt"
	"net/http"
	"os"
)

func ExampleLookupRequest() {
	uri := "http://localhost:8080/base-path/resource?first=false"
	req, _ := http.NewRequest("", uri, nil)

	name := ""

	s, err := LookupUrl("", nil)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = "test"
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	//s, err = LookupUrl("method", req.URL)
	//fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = SchemeName
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = HostName
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = PathName
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = QueryName
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	name = "$RUNTIME"
	os.Setenv("RUNTIME", "DEV")
	s, err = LookupUrl(name, req.URL)
	fmt.Printf("test: LookupUrl(%v) -> [err:%v] [%v]\n", name, err, s)

	//Output:
	//test: LookupUrl() -> [err:invalid argument: Url is nil] []
	//test: LookupUrl(test) -> [err:invalid argument : LookupEnv() template variable is invalid: test] []
	//test: LookupUrl(scheme) -> [err:<nil>] [http]
	//test: LookupUrl(host) -> [err:<nil>] [localhost:8080]
	//test: LookupUrl(path) -> [err:<nil>] [/base-path/resource]
	//test: LookupUrl(query) -> [err:<nil>] [first=false]
	//test: LookupUrl($RUNTIME) -> [err:<nil>] [DEV]

}
